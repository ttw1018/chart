import axios from "axios";

export async function getData(dataset: string, ratios: string[]) {
  const res = await axios({
    url: "/api/result",
    method: "get",
    params: {
      dataset: dataset,
      ratio: ratios,
    },
  });
  return res.data;
}

export async function getError(
  dataset: string,
  ratio: string[],
  n_clusters: string[]
) {
  const res = await axios({
    url: "/api/error",
    method: "get",
    params: {
      dataset: dataset,
      ratio: ratio,
      n_clusters: n_clusters,
    },
  });
  return res.data;
}

export async function getInfo() {
  const res = await axios({
    url: "/api/info",
    method: "get",
  });
  return res.data;
}

export async function getCluster(
  dataset: string,
  metric: string,
  feature: string
) {
  const res = await axios({
    url: "api/cluster",
    method: "get",
    params: {
      dataset: dataset,
      metric: metric,
      feature: feature,
    },
  });
  return res.data;
}
