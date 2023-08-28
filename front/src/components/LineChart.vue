<script setup lang="ts">
import { ref, onMounted } from "vue";
import * as echarts from "echarts";
import { getData } from "../utils/getData.ts";

const props = defineProps(["trainType"]);

const datasets =
  props.trainType == "few-shot"
    ? ["crosswoz", "sgd", "multiwoz2.2"]
    : ["multiwoz2.1", "sgd"];

const chartInfos =
  props.trainType == "few-shot" ? ["data ratio", "n_clusters"] : ["n_clusters"];

var chartInfo = ref(chartInfos[0]);

var dataset = ref(datasets[0]);

var chartDom;
var mychart: any;

var option = {
  title: {
    text: dataset.value + "-" + chartInfo.value,
  },
  toolbox: {
    show: true,
    feature: {
      saveAsImage: {},
    },
  },
  legend: {
    data: [] as string[],
  },
  xAxis: {
    type: "category",
    boundaryGap: false,
    data: [],
  },
  yAxis: {
    type: "value",
    scale: true,
    axisPointer: {
      snap: true,
    },
  },
  series: [] as { [key: string]: string }[],
};

onMounted(() => {
  chartDom = document.getElementById("linechart")!;
  mychart = echarts.init(chartDom);
  getChartInfo();
});

async function getChartInfo() {
  var infos;
  if (props.trainType == "zero-shot") {
    infos = ["zero-shot"];
  } else {
    infos = ["1", "5", "10", "100"];
  }
  var tmp = await getData(dataset.value, infos);

  var data = tmp["result"];

  var groupedData: { [key: string]: any } = {};

  if (chartInfo.value == "data ratio") {
    groupedData = data.reduce(
      (result: { [key: string]: any }[], item: { [key: string]: any }) => {
        const dataRatio = item["data ratio"];
        if (!result[dataRatio]) {
          result[dataRatio] = [];
        }
        result[dataRatio].push(item);
        return result;
      },
      {} as { [key: string]: any },
    );
  } else if (chartInfo.value == "n_clusters") {
    groupedData = data.reduce(
      (result: { [key: string]: any }[], item: { [key: string]: any }) => {
        const { n_clusters } = item;
        if (!result[n_clusters]) {
          result[n_clusters] = [];
        }
        result[n_clusters].push(item);
        return result;
      },
      {} as { [key: string]: any },
    );
  }

  option["legend"].data = Object.keys(groupedData);

  var series = [] as { [key: string]: any }[];

  Object.keys(groupedData).forEach((key: string) => {
    var r = {
      name: key,
      type: "line",
      data: [] as string[],
      label: {
        show: true,
      },
    } as { [key: string]: any };

    groupedData[key].forEach((item: { [key: string]: any }) => {
      r["data"].push(item["turn acc"]);
    });
    series.push(r);
  });

  option.series = series;
  console.log(option);

  mychart.setOption(option);
}
</script>

<template>
  <el-select
    v-model="dataset"
    filterable
    remote
    :remote-method="getChartInfo"
    size="large"
  >
    <el-option
      v-for="item in datasets"
      :key="item"
      :value="item"
      :label="item"
    />
  </el-select>

  <el-select
    v-model="chartInfo"
    filterable
    remote
    :remote-method="getChartInfo"
    size="large"
  >
    <el-option
      v-for="item in chartInfos"
      :key="item"
      :value="item"
      :label="item"
    />
  </el-select>

  <div
    id="linechart"
    style="width: 1000px; height: 800px; background-color: white"
  ></div>
</template>
