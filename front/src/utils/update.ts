import axios from "axios";
import { ElNotification } from "element-plus";

export async function update() {
  await axios.get("/api/update").then((res) => {
    let result = res.data.status
    ElNotification(
      result == "success"
        ? {
          title: "Success",
          message: "Successfully updata all data",
          type: "success",
        }
        : {
          title: "Fail",
          "message": "Fail to update data",
          type: "error",
        },
    );
  });
}
