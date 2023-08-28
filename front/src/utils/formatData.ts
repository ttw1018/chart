export function normalizeFloat(
    datas: { [key: string]: string }[],
    columns: string[],
) {
    for (var data of datas) {
        for (const column of columns) {
            data[column] = (parseFloat(data[column]) * 100).toFixed(2);
        }
    }
    return datas;
}

export function normalizeFloat1(
    datas: { [key: string]: string }[],
    columns: string[],
) {
    for (var data of datas) {
        for (const column of columns) {
            data[column] = (parseFloat(data[column])).toFixed(2);
        }
    }
    return datas;
}

export function normalizeInt(
    datas: { [key: string]: string }[],
    columns: string[],
) {
    for (var data of datas) {
        for (const column of columns) {
            data[column] = parseInt(data[column]).toString();
        }
    }
    return datas;
}