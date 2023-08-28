<script lang="ts">
import * as echarts from "echarts";
import { defineComponent } from "vue";
import data from "../data/exist_turn.json";

export default defineComponent({
  data() {
    return {
      tableData: [] as any[],
      dataset: "crosswoz",
      datas: data as { [key: string]: any },
      datasets: Object.keys(data),
      option: {
        title: {
          text: "",
          subtext: "value split ratio",
          left: "center",
        },
        tooltip: {
          trigger: "item",
        },
        toolbox: {
          show: true,
          feature: {
            saveAsImage: {},
          },
        },
        legend: {
          orient: "vertical",
          left: "left",
        },
        series: {
          name: "Access From",
          type: "pie",
          radius: "70%",
          data: [] as any[],
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: "rgba(0, 0, 0, 0.5)",
            },
          },
        },
      },
    };
  },
  methods: {
    draw() {
      this.setoption(this.dataset);
      var chartDom = document.getElementById("pie")!;
      var myChart = echarts.init(chartDom);
      myChart.setOption(this.option);
    },
    setoption(ds: string) {
      let data = this.datas[ds];
      var mydata: object[] = [];
      var sum: number = 0;
      for (var k in data) {
        mydata.push({
          value: data[k],
          name: k,
        });
        sum += data[k];
      }
      this.option.series.data = mydata;
      var tmp: number = 0;
      mydata.forEach((element: any) => {
        tmp += (element["value"] / sum) * 100;
        var ratio: string = String((element["value"] / sum) * 100);
        element["ratio"] = parseFloat(ratio).toFixed(2);
        element["sum"] = parseFloat(String(tmp)).toFixed(2);
      });
      this.tableData = mydata;
    },
  },
  mounted() {
    this.draw();
  },
});
</script>

<template>
  <div>
    <el-select
      v-model="dataset"
      filterable
      remote
      :remote-method="draw"
      placeholder="crosswoz"
      size="large"
    >
      <el-option
        v-for="item in datasets"
        :key="item"
        :value="item"
        :label="item"
      />
    </el-select>
  </div>

  <div id="pie" style="width: 600px; height: 600px; float: left"></div>

  <div>
    <el-table
      :data="tableData"
      style="width: 600px; float: left; margin-left: 20px"
    >
      <el-table-column prop="name" label="exist turn" width="120" />
      <el-table-column prop="value" label="count" width="120" />
      <el-table-column prop="ratio" label="ratio" width="120" />
      <el-table-column prop="sum" label="sum" width="120" />
    </el-table>
  </div>
</template>

<style></style>
