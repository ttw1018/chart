export const getCombineInfo = (
  data: { [key: string]: any }[],
  combineColumn: string[],
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
