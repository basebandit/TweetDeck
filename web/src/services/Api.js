import axios from "axios"

export default (token) => {
  if (token) {
    const request = axios.create({
      baseURL: `http://localhost:8880`,
      headers: {
        Authorization: "Bearer " + token,
        "Content-Type": "application/json",
      },
      timeout: 60000,
    })

    return request
  }

  return axios.create({
    baseURL: `http://localhost:8880`,
    timeout: 60000,
  })
}