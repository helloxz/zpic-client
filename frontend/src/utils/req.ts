import axios from 'axios';

const req = axios.create({
  timeout: 60000,
});

// 请求拦截器：每次请求时动态读取 base_url 和 token
req.interceptors.request.use(
  (config) => {
    const base_url = localStorage.getItem('base_url') || "https://www.imgurl.org";
    config.baseURL = base_url;

    const token = localStorage.getItem('token');
    if (token) {
      config.headers['Authorization'] = "Bearer " + token;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export function toForm(data:any) {
  const formData = new FormData();
  for (const key in data) {
      if (data.hasOwnProperty(key)) {
          formData.append(key, data[key]);
      }
  }
  return formData;
}

export default req;