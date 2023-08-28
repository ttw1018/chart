<script setup lang="ts">
import Table from "./components/table.vue";
import ErrorAnalysis from "./components/ErrorAnalysis.vue";
import Analysis from "./components/Analysis.vue";
import Cluster from "./components/cluster.vue";
import { onMounted, ref } from "vue";

console.log(localStorage.getItem("activateName"));
let activeName = ref<string>("few-shot");

var dd = ref(["a", "b", "c"]);
</script>

<template>
  <el-tabs v-model="activeName" class="demo-tabs">
    <el-tab-pane label="analysis" name="analysis">
      <Analysis />
    </el-tab-pane>
    <el-tab-pane label="cluster" name="cluster">
      <Cluster />
    </el-tab-pane>
    <el-tab-pane label="few-shot result" name="few-shot">
      <Table trainType="few-shot" :dddd="dd"></Table>
    </el-tab-pane>
    <el-tab-pane label="zero-shot result" name="zero-shot">
      <Table trainType="zero-shot"></Table>
    </el-tab-pane>
    <el-tab-pane label="error-analysis" name="error-analysis">
      <ErrorAnalysis></ErrorAnalysis>
    </el-tab-pane>
  </el-tabs>
</template>

<style>
.demo-tabs {
  margin-left: 8px;

  .el-tabs__content {
    padding: 12px;
    color: #6b778c;
    font-size: 24px;
  }
}
</style>
