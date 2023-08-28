<template>
    <h1>
        slotname distribution
    </h1>
    <el-table :data="all_info">
        <el-table-column v-for="item in columns" :prop="item" :label="item == 'index' ? 'dataset' : item" />
    </el-table>
</template>


<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { getInfo } from "../utils/getData"
import { normalizeFloat1, normalizeInt } from "../utils/formatData";

var all_info = ref<Array<{ [key: string]: any }>>([])

var columns: Array<string> = ["index", "count", "25%", "50%", "75%", "min", "max", "mean"]

const getdata = async () => {
    var tmp = await getInfo()
    all_info.value = normalizeFloat1(tmp["info"], ["mean", "std", "25%", "50%", "75%"])
    all_info.value = normalizeInt(all_info.value, ["count", "max", "min"])
}

onMounted(() => {
    getdata()
})

</script>