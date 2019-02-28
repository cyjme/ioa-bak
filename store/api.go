//generate by gen
package store

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
)

type Api struct {
	Id         string   `json:"id"`
	ApiGroupId string   `json:"apiGroupId"`
	Name       string   `json:"name" binding:"required"`
	Tags       []string `json:"tags"`
	Describe   string   `json:"describe"`
	Path       string   `json:"path" binding:"required"`
	Methods    []string `json:"methods" binding:"required"`
	Status     string   `json:"status"`

	Targets []Target `json:"targets"`
	Plugins string   `json:"plugins"`

	Policies     []string `json:"policies"`
	PoliciesData []Policy `json:"policiesData"`
}

type Target struct {
	Scheme string `json:"scheme"`
	Method string `json:"method"`
	Host   string `json:"host"`
	Port   string `json:"port"`
	Path   string `json:"path"`
}

func (api *Api) Put() error {
	if api.Targets == nil {
		api.Targets = make([]Target, 0)
	}
	apiByte, err := json.Marshal(api)
	if err != nil {
		log.Error(ERR_STORE_CRUD_API, err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	_, err = clientv3.NewKV(client).Put(ctx, apiPrefix+api.Id, string(apiByte))
	if err != nil {
		log.Error(ERR_STORE_CRUD_API, err)
		return err
	}
	cancel()

	return err
}

func (api *Api) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	_, err := clientv3.NewKV(client).Delete(ctx, apiPrefix+api.Id)
	if err != nil {
		log.Error(ERR_STORE_CRUD_API, err)
	}
	cancel()

	return err
}

func (api *Api) List() ([]Api, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	res, err := clientv3.NewKV(client).Get(ctx, apiPrefix, clientv3.WithPrefix())
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

		for _, policyId := range api.Policies {
			policy := &Policy{
				Id: policyId,
			}
			policy, err := policy.Get()
			if err != nil {
				log.Error(ERR_STORE_CRUD_API, err)
				continue
			}
			api.PoliciesData = append(api.PoliciesData, *policy)
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
	res, err := clientv3.NewKV(client).Get(ctx, apiPrefix+api.Id)
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
	responseWatchChannel := client.Watch(context.Background(), apiPrefix, clientv3.WithPrefix())

	for wresp := range responseWatchChannel {
		for _, ev := range wresp.Events {
			log.Info("watch api data change", ev.Type, ev.Kv.Key, ev.Kv.Value)

			callback()
		}
	}
}
