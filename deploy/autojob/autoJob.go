package autojob

import (
	"log"
	"mail/deploy/context"
	"mail/mongo/model"
	"time"

	uuid "github.com/satori/go.uuid"
)

type AutoJob struct {
}

func (j AutoJob) Do() {
	ctx := context.ContextSingle

	for {
		tick := 10
		deploy, err := ctx.MongoClient.DeployRepository().Find()
		if err != nil {
			log.Println("auto job err:", err)
			time.Sleep(time.Duration(int(time.Second) * tick))
			continue
		}
		if deploy == nil {
			log.Println("deploy is nil skip")
			time.Sleep(time.Duration(int(time.Second) * tick))
			continue
		}
		time.Sleep(time.Duration(int(time.Second) * tick))
		tick = deploy.RefreshTick
		if tick <= 0 {
			tick = 10
		}
		list, err := ctx.SDK.Containers()
		if err != nil {
			log.Println("auto job get container err:", err)
			continue
		}
		if len(list) == 0 {
			ct := &model.Container{
				Id:   uuid.NewV4().String(),
				Ip:   ctx.Config.HostIp,
				List: []model.ContainerInfo{},
			}
			log.Println("auto job get container is nil")
			err = ctx.MongoClient.ContainerRepository().UpsertByField("ip", ctx.Config.HostIp, ct)
			if err != nil {
				log.Println(err)
			}
			continue
		}

		cList := make([]model.ContainerInfo, 0)
		for _, v := range list {
			macvlan := v.NetworkSettings.Networks["macvlan"]
			var ip string
			if macvlan != nil {
				ip = macvlan.IPAddress
			}
			c := model.ContainerInfo{
				Name:        v.Names[0],
				Id:          v.ID,
				Ip:          ip,
				Status:      v.Status,
				CreatedTime: time.Now().Unix(),
			}
			cList = append(cList, c)
		}
		ct := &model.Container{
			Id:   uuid.NewV4().String(),
			Ip:   ctx.Config.HostIp,
			List: cList,
		}
		dbContainer, err := ctx.MongoClient.ContainerRepository().FindByField("ip", ctx.Config.HostIp)
		if err != nil {
			log.Println("auto job get db container err:", err)
			continue
		}
		if dbContainer != nil && len(dbContainer.List) != 0 {
			oldList := dbContainer.List
			for i := 0; i < len(cList); i++ {
				for _, v := range oldList {
					if v.Id == cList[i].Id {
						cList[i].PptpIp = v.PptpIp
					}
				}
			}
		}
		err = ctx.MongoClient.ContainerRepository().UpsertByField("ip", ctx.Config.HostIp, ct)
		if err != nil {
			log.Println("auto job get save container err:", err)
			continue
		}
		log.Println("auto job process success")
	}

}
