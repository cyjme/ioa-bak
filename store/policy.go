package store

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
)

type Policy struct {
	Id      string `json:"id"`
	Name    string `json:"name" binding:"required"`
	Plugins string `json:"plugins"`
	Status  string `json:"status"`
}

func (policy *Policy) Put() error {
	policyByte, err := json.Marshal(policy)
	if err != nil {
		log.Error(ERR_STORE_CRUD_API, err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	_, err = clientv3.NewKV(client).Put(ctx, policyPrefix+policy.Id, string(policyByte))
	if err != nil {
		log.Error(ERR_STORE_CRUD_API, err)
		return err
	}
	cancel()

	return err
}

func (policy *Policy) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	_, err := clientv3.NewKV(client).Delete(ctx, policyPrefix+policy.Id)
	if err != nil {
		log.Error(ERR_STORE_CRUD_API, err)
	}
	cancel()

	return err
}

func (policy *Policy) List() ([]Policy, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	res, err := clientv3.NewKV(client).Get(ctx, policyPrefix, clientv3.WithPrefix())
	cancel()
	if err != nil {
		return nil, 0, err
	}
	policies := make([]Policy, 0)
	for _, kv := range res.Kvs {
		policy := Policy{}
		err := json.Unmarshal(kv.Value, &policy)
		if err != nil {
			return nil, 0, err
		}
		policies = append(policies, policy)
	}

	if err != nil {
		return nil, 0, err
	}
	return policies, len(policies), nil
}

func (policy *Policy) Get() (*Policy, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	res, err := clientv3.NewKV(client).Get(ctx, policyPrefix+policy.Id)
	cancel()
	if err != nil {
		return nil, err
	}

	if len(res.Kvs) == 0 {
		return nil, ERR_STORE_POLICY_NOT_FOUND
	}

	err = json.Unmarshal(res.Kvs[0].Value, policy)
	if err != nil {
		return nil, err
	}
	return policy, err
}

func (policy *Policy) Watch(callback func()) {
	responseWatchChannel := client.Watch(context.Background(), policyPrefix, clientv3.WithPrefix())

	for wresp := range responseWatchChannel {
		for _, ev := range wresp.Events {
			log.Info("watch policy data change", ev.Type, ev.Kv.Key, ev.Kv.Value)

			callback()
		}
	}
}
