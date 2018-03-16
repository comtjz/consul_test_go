package main

import (
	consulapi "github.com/hashicorp/consul/api"
	"log"
	"strings"
	"flag"
)

func main() {
	var consul_addr, monitor_addr, service_name, service_ip string
	var service_port int
	flag.StringVar(&consul_addr, "consul_addr", "localhost:8500", "host:port of the consul")
	flag.StringVar(&monitor_addr, "monitor_addr", "127.0.0.1:54321", "host:port of the service monitor")
	flag.StringVar(&service_name, "service_name", "worker1", "name of the service")
	flag.StringVar(&service_ip, "service_ip", "127.0.0.1", "service serve ip")
	flag.IntVar(&service_port, "service_port", 10086, "service serve port")
	flag.Parse()

	//DoRegistService(consul_addr, consul_addr, service_name, service_ip, service_port)
	//DiscoverService(consul_addr, true, "")
	//DoDegisterService(consul_addr, "")
}

func RegisterService(consul_addr string)

func DoDegisterService(consul_addr string, service_id string) {
	consulConf := consulapi.DefaultConfig()
	consulConf.Address = consul_addr
	client, err := consulapi.NewClient(consulConf)
	if err != nil {
		log.Fatal(err.Error())
	}

	client.Agent().ServiceDeregister("worker-127.0.0.1")
}

/**
 *
 */

/**
 * 连接指定consul，然后注册服务
 * consul_addr: 指定的consul的地址
 * monitor_addr:
 */
func DoRegistService(consul_addr string, monitor_addr string, service_name string, ip string, port int) {
	//my_service_id := service_name + "-" + ip
	my_service_id := "worker-"+ip
	var tags []string
	service := &consulapi.AgentServiceRegistration{
		ID: my_service_id,
		Name: service_name,
		Port: port,
		Address: ip,
		Tags: tags,
		/*
		Check: &consulapi.AgentServiceCheck {
			HTTP: "http://" + monitor_addr + "/status",
			Interval: "5s",
			Timeout: "1s",
		},
		*/
	}

	consulConf := consulapi.DefaultConfig()
	consulConf.Address = consul_addr
	client, err := consulapi.NewClient(consulConf)
	if err != nil {
		log.Fatal("1 " + err.Error())
	}

	if err := client.Agent().ServiceRegister(service); err != nil {
		log.Fatal("2 " + err.Error())
	}
	log.Printf("Registered service %q in consul with tags %q", service_name, strings.Join(tags, ","))
}

func DiscoverService(addr string, healthOnly bool, service_name string) {
	consulConf := consulapi.DefaultConfig()
	consulConf.Address = addr
	client, err := consulapi.NewClient(consulConf)
	if err != nil {
		log.Fatal(err.Error())
	}

	services, _, err := client.Catalog().Services(&consulapi.QueryOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}

	for name, _ := range services {
		log.Printf("service_name = %v\n", name)
		servicesData, _, err := client.Health().Service(name, "", healthOnly, &consulapi.QueryOptions{})
		if err != nil {
			log.Fatal(err.Error())
		}

		for _, entry := range servicesData {
			log.Printf("ID = %v\n", entry.Service.ID)
			log.Printf("Name = %v\n", entry.Service.Service)
		}
	}
}
