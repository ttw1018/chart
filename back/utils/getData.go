package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strings"
)

func readCSV(path string) []map[string]interface{} {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("open file failed: ", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	columns, err := reader.Read()

	if err != nil {
		fmt.Println("read file failed: ", err)
	}

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("read file failed: ", err)
	}

	var datas []map[string]interface{}
	for _, record := range records {
		data := make(map[string]interface{})
		for i := range columns {
			data[columns[i]] = record[i]
		}
		datas = append(datas, data)
	}
	return datas
}

func GetData(dataset string, ratios []string) []map[string]interface{} {
	var datas []map[string]interface{}
	for _, ratio := range ratios {
		datas = append(datas, getResult(dataset, ratio)...)
	}
	return datas
}

func getResult(dataset, ratio string) []map[string]interface{} {
	datapath := path.Join("data", ratio, dataset+".csv")

	file, err := os.Open(datapath)
	if err != nil {
		fmt.Println("open file failed: ", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	columns, err := reader.Read()

	if err != nil {
		fmt.Println("read file failed: ", err)
	}
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("read file failed: ", err)
	}

	var datas []map[string]interface{}

	var children []map[string]interface{}
	var dev []map[string]interface{}

	for _, record := range records {
		data := make(map[string]interface{})
		for i := range columns {
			data[columns[i]] = record[i]
		}
		if len(data["checkpoint"].(string)) > 0 {
			if ratio == "zero-shot" {
				data["domain"] = strings.Split(data["checkpoint"].(string), "/")[2]
			}
			if data["cluster"] == "all" {
				data["children"] = children
				data["dev"] = dev
				datas = append(datas, data)
				children = []map[string]interface{}{}
				dev = []map[string]interface{}{}
			} else {
				children = append(children, data)
			}
		} else {
			dev = append(dev, data)
		}
	}
	return datas
}

func GetError(dataset string, ratios, n_clusters []string) []map[string]interface{} {
	var data = make(map[string]map[string]map[string]interface{})

	for _, ratio := range ratios {
		for _, cluster := range n_clusters {
			datapath := path.Join("data", "analysis", fmt.Sprintf("%s-%s-%scluster-error-distribution.csv", dataset, ratio, cluster))

			file, err := os.Open(datapath)
			if err != nil {
				fmt.Println("open file failed: ", err)
			}
			defer file.Close()
			reader := csv.NewReader(file)

			columns, err := reader.Read()

			if err != nil {
				fmt.Println("read file failed: ", err)
			}
			records, err := reader.ReadAll()

			if err != nil {
				fmt.Println("read file failed: ", err)
			}

			for _, record := range records {
				d := make(map[string]string)
				for i := range columns {
					d[columns[i]] = record[i]
				}
				domain := d["domain"]
				slotname := d["slotname"]
				if _, ok := data[domain]; !ok {
					data[domain] = make(map[string]map[string]interface{})
				}
				if _, ok := data[domain][slotname]; !ok {
					data[domain][slotname] = make(map[string]interface{})
				}
				data[domain][slotname][cluster+"-domain_acc"] = d["domain_acc"]
				data[domain][slotname][cluster+"-slotname_acc"] = d["slotname_acc"]
				data[domain][slotname]["cluster"] = cluster
			}
		}
	}

	var datas []map[string]interface{}

	for domain, domainMap := range data {
		for slotname, slotnameMap := range domainMap {
			slotnameMap["domain"] = domain
			slotnameMap["slotname"] = slotname
			datas = append(datas, slotnameMap)
		}
	}
	return datas
}

func GetSlotnameDistribution() []map[string]interface{} {
	path := "./data/info/slotname_distribution.csv"

	data := readCSV(path)

	return data
}

func GetClusterInfo(dataset, metric, feature string) []map[string]interface{} {
	path := path.Join("data", "cluster", metric+"-"+feature+"-"+dataset+".csv")
	return readCSV(path)
}
