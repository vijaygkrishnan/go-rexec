package rexec

import ( 
         "testing"
         "io/ioutil"
       )

func TestParseConfigFile(t *testing.T) {
    want := []string { "swmp1-spine1.domain.local", "swmp1-spine2.domain.local" }
    err := ioutil.WriteFile("test_rexec.json", []byte(SampleInputJson), 0644)
    if err != nil {
        t.Fatal(err)
    }
    var got []string
    got = GetNodesFromCfgFile("swmp1-spines", "test_rexec.json")
    for i,_ := range(want) {
        if want[i] != got[i] {
            t.Fatal(err)
        }
    }
}
