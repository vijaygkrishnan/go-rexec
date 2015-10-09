package rexec

import gexpect "github.com/ThomasRooney/gexpect"
import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"time"
)

var SampleInputJson string = `
{  "Groups":
    [
        {
            "GroupName" : "swmp1-spines",
            "Nodes" : [
                { "NodeIP" : "swmp1-spine1.domain.local" },
                { "NodeIP" : "swmp1-spine2.domain.local" }
            ]
        },
        {
            "GroupName" : "swmp1-leafs",
            "Nodes" : [
                { "NodeIP" : "swmp1-leaf1.domain.local" },
                { "NodeIP" : "192.168.10.100" }
            ]
        },
        {
            "GroupName" : "swmp1",
            "Nodes" : [
                { "NodeIP" : "swmp1-spine1.domain.local" },
                { "NodeIP" : "swmp1-spine2.domain.local" },
                { "NodeIP" : "swmp1-leaf1.domain.local" },
                { "NodeIP" : "192.168.10.100" }
            ]
        }
    ]
}
`

type Config struct {
	Groups []Group
}
type Group struct {
	GroupName string
	Nodes     []Node
}
type Node struct {
	NodeIP string
}

func ParseConfigFile(cfgFilename string) Config {

	file, err := ioutil.ReadFile(cfgFilename)
	if err != nil {
		log.Fatal(err)
	}
	var ConfigFile Config
	err = json.Unmarshal(file, &ConfigFile)
	if err != nil {
		log.Fatalf("Error parsing %s: %v", cfgFilename, err)
	}
	return ConfigFile
}

func GetNodesFromCfgFile(GrpName string, cfgFilename string) []string {
	var NodeList []string
	cfg := ParseConfigFile(cfgFilename)
	for _, grp := range cfg.Groups {
		if grp.GroupName != GrpName {
			continue
		}
		for _, Node := range grp.Nodes {
			NodeList = append(NodeList, Node.NodeIP)
		}
	}
	return NodeList
}

func RemoteExecute(user string, password string, remoteIp string, cmd string) (string, error) {
	str := fmt.Sprintf("ssh %s@%s", user, remoteIp)
	ssh, err := gexpect.Spawn(str)
	if err != nil {
		fmt.Println("Spawn", str)
		fmt.Println("Could not connect to ", remoteIp, err)
		return "", err
	}
	ssh.Expect("sword")
	ssh.SendLine(password)
	err = ssh.ExpectTimeout("#", 10*time.Second)
	if err != nil {
		fmt.Println("Spawn", str)
		fmt.Println("password:", password)
		fmt.Println("Login failed to", remoteIp, "(", err, ")", "Please check password")
		return "", err
	}
	ssh.Capture()
	ssh.SendLine(cmd)
	ssh.Expect("#")
	bytes := ssh.Collect()
	ssh.SendLine("exit")
	ssh.Expect("#")
	outStr := fmt.Sprintf("Executed on %s: ", remoteIp) + string(bytes[:])
	return outStr, err
}

func TestbedRemoteExec(user string, passwd string, nodes []string, cmd string) []string {
	resultChan := make(chan string, len(nodes))
	for _, nodeName := range nodes {
		go func(nodeName string) {
			str, err := RemoteExecute(user, passwd, nodeName, cmd)
			if err != nil {
				fmt.Sprintf(str, err)
			}
			resultChan <- str
		}(nodeName)
	}
	resultList := make([]string, len(nodes))
	for i, _ := range nodes {
		resultList[i] = <-resultChan
	}
	return resultList
}

func TestbedRemoteExecSorted(user string, passwd string, nodes []string, cmd string) []string {
	resultList := TestbedRemoteExec(user, passwd, nodes, cmd)
	sort.Sort(sort.StringSlice(resultList))
	return resultList
}

func ParseCommandLineArgs() (user string, passwd string, group string, sort bool, cmd string, json bool) {
	p_user := flag.String("user", "root", "username of remote device")
	p_passwd := flag.String("passwd", "root", "password of remote device")
	p_group := flag.String("group", "swmp1", "group name from testbed.json")
	p_sort := flag.Bool("sort", false, "sort the output")
	p_cmd := flag.String("cmd", "", "Command to exec on remote device")
	p_json := flag.Bool("json", false, "Print a sample json file (testbed.json)")

	flag.Parse()
	if flag.NArg() != 0 {
		flag.Usage()
		log.Fatalf("Invalid argument: %s (Num Arg = %v)", flag.Args(), flag.NArg())
	}

	println("user:", *p_user)
	println("passwd:", *p_passwd)
	println("group:", *p_group)
	println("sort:", *p_sort)
	println("cmd:", *p_cmd)
	println("json:", *p_json)
	return *p_user, *p_passwd, *p_group, *p_sort, *p_cmd, *p_json
}

func Main() {
	var outputStrList []string
	user, passwd, group, sort, cmd, json := ParseCommandLineArgs()
	if json {
		println(SampleInputJson)
		return
	}
	nodelist := GetNodesFromCfgFile(group, "testbed.json")
	if sort {
		outputStrList = TestbedRemoteExecSorted(user, passwd, nodelist, cmd)
	} else {
		outputStrList = TestbedRemoteExec(user, passwd, nodelist, cmd)
	}
	for i, _ := range outputStrList {
		println(outputStrList[i])
	}
}
