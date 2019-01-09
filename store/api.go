//generate by gen
package store

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"ioa/httpServer/pkg/util"
	"log"
)

const prefix = "apis/"

type Api struct {
	Id         string `json:"id"`
	ApiGroupId string `json:"apiGroupId"`
	Name       string `json:"name"`
	Describe   string `json:"describe"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	Status     string `json:"status"`

	Targets []Target `json:"targets"`
	Plugins string   `json:"plugins"`
}

type Target struct {
	Scheme string `json:"scheme"`
	Method string `json:"method"`
	Host   string `json:"host"`
	Port   string `json:"port"`
	Path   string `json:"path"`
}

func (api *Api) Put() error {
	api.Id = util.GetRandomString(10)
	apiByte, err := json.Marshal(api)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	_, err = clientv3.NewKV(client).Put(ctx, prefix+api.Id, string(apiByte))
	cancel()

	return err
}

func (api *Api) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	_, err := clientv3.NewKV(client).Delete(ctx, prefix+api.Id)
	cancel()

	return err
}

func (api *Api) List() ([]Api, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	res, err := clientv3.NewKV(client).Get(ctx, prefix, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, 0, err
	}
	apis := make([]Api, 0)
	for _, kv := range res.Kvs {
		api := Api{}
		err := json.Unmarshal(kv.Value, &api)
		if err != nil {
			return nil, 0, err
		}
		apis = append(apis, api)
	}

	if err != nil {
		return nil, 0, err
	}
	return apis, len(apis), nil
}

func (api *Api) Get() (*Api, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	res, err := clientv3.NewKV(client).Get(ctx, prefix+api.Id)
	cancel()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res.Kvs[0].Value, api)
	if err != nil {
		return nil, err
	}
	return api, err
}

func (api *Api) Watch(callback func()) {
	responseWatchChannel := client.Watch(context.Background(), prefix, clientv3.WithPrefix())

	for wresp := range responseWatchChannel {
		for _, ev := range wresp.Events {
			log.Println("watch api change")
			log.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

			callback()
		}
	}
}
