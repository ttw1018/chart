<template>
  <div>
    <el-button type="primary" plain @click="update">update data</el-button>
    <el-cascader v-model="option" :props="props" @change="updatePage">
    </el-cascader>
    <el-select v-model="sortKey" filterable remote :remote-method="sortData">
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
      style="width: 100%; border: 1px solid black"
      :span-method="mergeRow"
      :cell-style="{ borderColor: 'black' }"
      :header-cell-style="{ borderColor: 'black' }"
      border
    >
      <el-table-column
        v-for="column in show_column"
        :label="column"
        :prop="column"
      >
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { getError } from "../utils/getData";
import { CascaderProps, TableColumnCtx } from "element-plus";
import { normalizeFloat } from "../utils/formatData";
import { getCombineInfo } from "../utils/mergeRow";
import { update } from "../utils/update.ts";

interface rowMergeProps {
  row: { [key: string]: any };
  column: TableColumnCtx<{ [key: string]: any }>;
  rowIndex: number;
  columnIndex: number;
}

const mergeRow = ({ row, column, rowIndex, columnIndex }: rowMergeProps) => {
  const columnName = column.property;
  if (combineInfo.value.hasOwnProperty(columnName)) {
    return [combineInfo.value[columnName][rowIndex], 1];
  }
};

const props: CascaderProps = {
  lazy: true,
  lazyLoad(node, resolve) {
    const { level } = node;
    if (level == 0) {
      var nodes = [
        {
          value: "few-shot",
          label: "few-shot",
          leaf: false,
        },
        {
          value: "zero-shot",
          label: "zero-shot",
          leaf: false,
        },
      ];
      resolve(nodes);
    } else if (level == 1) {
      var nodes = [] as { value: string; label: string; leaf: boolean }[];
      if (node.data?.value == "few-shot") {
        for (const ds of few_shot_datasets) {
          nodes.push({
            value: ds,
            label: ds,
            leaf: false,
          });
        }
      } else if (node.data?.value == "zero-shot") {
        for (const ds of zero_shot_datasets) {
          nodes.push({
            value: ds,
            label: ds,
            leaf: true,
          });
        }
      }
      resolve(nodes);
    } else if (level == 2) {
      var nodes = [] as { value: string; label: string; leaf: boolean }[];
      if (node.parent?.data?.value == "few-shot") {
        for (const ratio of few_shot_ratios) {
          nodes.push({
            value: ratio,
            label: ratio,
            leaf: true,
          });
        }
      } else if (node.parent?.data?.value == "zero-shot") {
        for (const ratio of zero_shot_ratio) {
          nodes.push({
            value: ratio,
            label: ratio,
            leaf: true,
          });
        }
      }
      resolve(nodes);
    }
  },
};

var show_column = ref([
  "cluster",
  "domain",
  "domain_acc",
  "slotname",
  "slotname_acc",
]);

const few_shot_datasets = ["crosswoz", "sgd", "multiwoz2.2"];
const zero_shot_datasets = ["sdg", "multiwoz2.1"];
const few_shot_ratios = ["1", "5", "10", "100"];
const zero_shot_ratio = ["attraction", "hotel", "restaurant", "taxi", "train"];
const n_clusters = ["1", "2", "3", "5"];

const combineColumn = [
  "domain",
  "slotname",
  "cluster",
  "1-domain_acc",
  "2-domain_acc",
  "3-domain_acc",
  "5-domain_acc",
];

var option = ref([] as string[]);
var data = ref([] as { [key: string]: any }[]);

var combineInfo = ref({} as { [key: string]: any });

async function updatePage() {
  if (option.value[0] == "few-shot") {
    var tmp = await getError(option.value[1], [option.value[2]], n_clusters);
  } else if (option.value[0] == "zero-shot") {
    var tmp = await getError(
      option.value[1],
      ["attraction", "hotel", "restaurant", "taxi", "train"],
      n_clusters
    );
  }
  data.value = tmp["error"];
  data.value = normalizeFloat(data.value, [
    "1-domain_acc",
    "2-domain_acc",
    "3-domain_acc",
    "5-domain_acc",
    "1-slotname_acc",
    "2-slotname_acc",
    "3-slotname_acc",
    "5-slotname_acc",
  ]);
  sortData();
}

const sortKeys = ["domain", "cluster"];
var sortKey = ref(sortKeys[0]);

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
  if (key1 == "domain") {
    show_column.value = [
      "domain",
      "1-domain_acc",
      "2-domain_acc",
      "3-domain_acc",
      "5-domain_acc",
      "slotname",
      "1-slotname_acc",
      "2-slotname_acc",
      "3-slotname_acc",
      "5-slotname_acc",
    ];
  } else {
    show_column.value = [
      "cluster",
      "domain",
      "slotname",
      "1-slotname_acc",
      "2-slotname_acc",
      "3-slotname_acc",
      "5-slotname_acc",
    ];
  }
  combineInfo.value = getCombineInfo(data.value, combineColumn);
}

onMounted(() => {
  option.value = ["few-shot", few_shot_datasets[0], few_shot_ratios[0]];
  updatePage();
});
</script>
