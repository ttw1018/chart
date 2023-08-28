<template>
  <div>
    <el-select v-model="dataset" filterable remote :remote-method="setOption">
      <el-option
        v-for="item in datasets"
        :key="item"
        :value="item"
        :label="item"
      />
    </el-select>
  </div>

  <div id="chart1" style="width: 800px; height: 500px; float: left"></div>
  <div id="chart2" style="width: 800px; height: 500px; float: left"></div>
  <div id="chart3" style="width: 800px; height: 500px; float: left"></div>
  <div id="chart4" style="width: 800px; height: 500px; float: left"></div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { getCluster } from "../utils/getData";
import * as echarts from "echarts";

const datasets = [
  "crosswoz",
  "sgd",
  "multiwoz2.2",
  "multiwoz2.1-hotel",
  "multiwoz2.1-train",
  "multiwoz2.1-attraction",
  "multiwoz2.1-restaurant",
  "multiwoz2.1-taxi",
  "multiwoz2.2-hotel",
  "multiwoz2.2-train",
  "multiwoz2.2-attraction",
  "multiwoz2.2-restaurant",
  "multiwoz2.2-taxi",
];
const features = ["transformer", "embedding"];
const metrics = ["SSE", "silhouette"];
var dataset = ref<string>(datasets[0]);
var metric = ref<string>(metrics[0]);
var feature = ref<string>(features[0]);

var data = ref<{ [key: string]: any }>({});
var chart1Dom: any;
var mychart1: echarts.EChartsType;

var chart2Dom: any;
var mychart2: echarts.EChartsType;

var chart3Dom: any;
var mychart3: echarts.EChartsType;

var chart4Dom: any;
var mychart4: echarts.EChartsType;

var option1 = ref({
  title: {
    text: "",
  },
  xAxis: {
    type: "category",
    data: [] as string[],
  },
  legend: {
    data: [] as string[],
  },
  yAxis: {
    type: "value",
  },
  series: [] as any[],
});

var option2 = ref({
  title: {
    text: "",
  },
  xAxis: {
    type: "category",
    data: [] as string[],
  },
  legend: {
    data: [] as string[],
  },
  yAxis: {
    type: "value",
  },
  series: [] as any[],
});
var option3 = ref({
  title: {
    text: "",
  },
  xAxis: {
    type: "category",
    data: [] as string[],
  },
  legend: {
    data: [] as string[],
  },
  yAxis: {
    type: "value",
  },
  series: [] as any[],
});
var option4 = ref({
  title: {
    text: "",
  },
  xAxis: {
    type: "category",
    data: [] as string[],
  },
  legend: {
    data: [] as string[],
  },
  yAxis: {
    type: "value",
  },
  series: [] as any[],
});

onMounted(() => {
  chart1Dom = document.getElementById("chart1");
  mychart1 = echarts.init(chart1Dom);

  chart2Dom = document.getElementById("chart2");
  mychart2 = echarts.init(chart2Dom);

  chart3Dom = document.getElementById("chart3");
  mychart3 = echarts.init(chart3Dom);

  chart4Dom = document.getElementById("chart4");
  mychart4 = echarts.init(chart4Dom);

  setOption();
});

async function setOption() {
  var da = dataset.value;
  for (const me of metrics) {
    for (const fe of features) {
      if (!data.value.hasOwnProperty(da)) {
        data.value[da] = {};
      }
      if (!data.value[da].hasOwnProperty(me)) {
        data.value[da][me] = {};
      }
      if (!data.value[da][me].hasOwnProperty(fe)) {
        var tmp = await getCluster(da, me, fe);
        data.value[da][me][fe] = tmp["info"];
      }
    }
  }

  const d1: any[] = data.value[da][metrics[0]][features[0]];
  option1.value["xAxis"]["data"] = d1.map((x) => x["n_clusters"]).slice(0, 10);
  option1.value["series"] = [
    {
      data: d1.map((x) => x["SSE"]).slice(0, 10),
      type: "line",
    },
  ];
  option1.value.title.text = da + "-" + metrics[0] + "-" + features[0];
  mychart1.setOption(option1.value);

  const d2: any[] = data.value[da][metrics[1]][features[0]];
  option2.value["xAxis"]["data"] = d2.map((x) => x["n_clusters"]).slice(0, 10);
  option2.value["series"] = [
    {
      data: d2.map((x) => x["silhouette"]).slice(0, 10),
      type: "line",
    },
  ];
  option2.value.title.text = da + "-" + metrics[1] + "-" + features[0];
  mychart2.setOption(option2.value);

  const d3: any[] = data.value[da][metrics[0]][features[1]];
  option3.value["xAxis"]["data"] = d3.map((x) => x["n_clusters"]).slice(0, 10);
  option3.value["series"] = [
    {
      data: d3.map((x) => x["SSE"]).slice(0, 10),
      type: "line",
    },
  ];
  option3.value.title.text = da + "-" + metrics[0] + "-" + features[1];
  mychart3.setOption(option3.value);

  const d4: any[] = data.value[da][metrics[1]][features[1]];
  option4.value["xAxis"]["data"] = d4.map((x) => x["n_clusters"]).slice(0, 10);
  option4.value["series"] = [
    {
      data: d4.map((x) => x["silhouette"]).slice(0, 10),
      type: "line",
    },
  ];
  option4.value.title.text = da + "-" + metrics[1] + "-" + features[1];
  mychart4.setOption(option4.value);
}
</script>
