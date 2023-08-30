<template>
  <div>
    <el-button type="primary" plain @click="update">update data</el-button>
    <el-select
      v-model="dataset"
      filterable
      remote
      :remote-method="updateData"
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

    <el-select
      v-model="sortKey"
      filterable
      remote
      :remote-method="sortData"
      placeholder="sortKeys[0]"
      size="large"
    >
      <el-option
        v-for="item in sortKeys"
        :key="item"
        :value="item"
        :label="item"
      />
    </el-select>
  </div>

  <div>
    <el-table
      :data="data"
      border
      stripe
      style="
        width: 100%;
        float: left;
        margin-left: 20px;
        margin-right: 20px;
        border: 1px solid black;
      "
      :cell-style="{ borderColor: 'black' }"
      :header-cell-style="{ borderColor: 'black' }"
      :span-method="mergeRow"
    >
      <el-table-column type="expand">
        <template #default="props">
          <h2>
            <p style="margin-left: 20px">detail information</p>
          </h2>
          <el-table :data="props.row.children" style="margin-left: 20px">
            <el-table-column prop="cluster" label="cluster" width="120" />
            <el-table-column prop="slot acc" label="slot acc" width="120" />
            <el-table-column prop="turn acc" label="turn acc" width="120" />
            <el-table-column
              prop="dialogue acc"
              label="dialogue acc"
              width="120"
            />
          </el-table>
          <div
            :id="dataset + '-cluster-' + num + '-' + props.row.index"
            v-for="num in Array.from(
              { length: props.row.n_clusters },
              (_, index) => index
            )"
            style="width: 800px; height: 500px"
          ></div>
        </template>
      </el-table-column>

      <el-table-column prop="time" label="time" width="200" />
      <el-table-column prop="history turn" label="history turn" width="120" />
      <el-table-column
        v-if="props.trainType == 'few-shot'"
        prop="data ratio"
        label="ratio"
        width="120"
      />
      <el-table-column
        v-if="props.trainType == 'zero-shot'"
        prop="domain"
        label="domain"
        width="120"
      />
      <el-table-column prop="feature" label="feature" width="120" />
      <el-table-column prop="n_clusters" label="n_clusters" width="120" />
      <el-table-column prop="slot acc" label="slot acc" width="120" />
      <el-table-column prop="turn acc" label="turn acc" width="120" />
      <el-table-column prop="dialogue acc" label="dialogue acc" width="120" />
      <el-table-column prop="checkpoint" label="checkpoint" />
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { update } from "../utils/update";
import type { TableColumnCtx } from "element-plus";
import { getData } from "../utils/getData";

const props = defineProps(["trainType"]);

const datasets =
  props.trainType == "few-shot"
    ? ["crosswoz", "multiwoz2.2", "sgd"]
    : ["multiwoz2.1", "multiwoz2.2", "sgd"];

const sortKeys =
  props.trainType == "few-shot"
    ? ["data ratio", "n_clusters"]
    : ["domain", "n_clusters"];

var sortKey = ref(sortKeys[0]);

var dataset = ref(datasets[0]);

let data = ref([] as { [key: string]: any }[]);

let result: { [key: string]: any } = reactive({});

let combineColumn = ["dataset", "history turn", "data ratio", "n_clusters"];
// let combineColumn = ["data ratio", "n_clusters", "domain"] as string[];

let combineInfo: any = {};

async function getdata(dataset: any) {
  if (!result.hasOwnProperty(dataset)) {
    result[dataset] = {};
  }
  var ratios: string[];
  if (props.trainType == "zero-shot") {
    ratios = ["zero-shot"];
  } else {
    ratios = ["1", "5", "10", "100"];
  }
  var tmp = await getData(dataset, ratios);

  result[dataset] = tmp["result"];

  data.value = result[dataset];

  sortData();
}

function updateData() {
  getdata(dataset.value);
}

onMounted(() => {
  updateData();
});

interface Result {
  time: string;
  history: number;
  "slot acc": number;
  "turn acc": number;
  "dialogue history": number;
  checkpoint: string;
}

interface rowMergeProps {
  row: Result;
  column: TableColumnCtx<Result>;
  rowIndex: number;
  columnIndex: number;
}

const getCombineInfo = (
  data: { [key: string]: any }[],
  combineColumn: string[]
) => {
  var combineResult: any = {};
  var count: any = {};
  for (let i of combineColumn) {
    combineResult[i] = [];
    count[i] = 1;
  }
  for (var i = data.length - 1; i >= 0; i--) {
    for (let v of combineColumn) {
      if (i == 0) {
        combineResult[v].push(count[v]);
      } else {
        if (data[i][v] == data[i - 1][v]) {
          count[v] = count[v] + 1;
          combineResult[v].push(0);
        } else {
          combineResult[v].push(count[v]);
          count[v] = 1;
        }
      }
    }
  }
  for (const i in combineResult) {
    combineResult[i] = combineResult[i].reverse();
  }
  return combineResult;
};

const mergeRow = ({ row, column, rowIndex, columnIndex }: rowMergeProps) => {
  let columnName = column.property;
  if (combineInfo.hasOwnProperty(columnName)) {
    return [combineInfo[columnName][rowIndex], 1];
  }
};

function sortData() {
  var key1: string, key2: string;
  if (sortKey.value == sortKeys[0]) {
    key1 = sortKeys[0];
    key2 = sortKeys[1];
  } else {
    key1 = sortKeys[1];
    key2 = sortKeys[0];
  }
  data.value = data.value.sort(
    (a: { [key: string]: any }, b: { [key: string]: any }) => {
      if (a[key1] == b[key1]) {
        if (key2 == "domain") {
          return a[key2].localeCompare(b[key2]);
        } else {
          return a[key2] - b[key2];
        }
      }
      if (key1 == "domain") {
        return a[key1].localeCompare(b[key1]);
      } else {
        return a[key1] - b[key1];
      }
    }
  );
  combineInfo = getCombineInfo(data.value, combineColumn);
}
</script>
