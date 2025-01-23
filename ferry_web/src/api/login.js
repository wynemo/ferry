import request from "@/utils/request";

// 获取验证码
export function getCodeImg() {
  return request({
    url: "/api/v1/getCaptcha",
    method: "get",
  });
}

// 获取验证码
export function getSMS(phone) {
  return request({
    url: "/api/v1/getSMS",
    method: "get",
    params: { phone: phone },
  });
}
