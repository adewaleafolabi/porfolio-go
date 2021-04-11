import axios from "axios";
import * as qs from "qs";
// import {cacheAdapterEnhancer, throttleAdapterEnhancer} from 'axios-extensions';

export const BASE_URL = process.env.VUE_APP_API_URL;
export const API_SERVICE = axios.create({
    baseURL: BASE_URL,
    headers: {},
    timeout: 20 * 1000,
    withCredentials: true,
    // adapter: throttleAdapterEnhancer(cacheAdapterEnhancer(axios.defaults.adapter))
})

export const getPortfolios = async ()=>{
    return getList("/portfolios")
}

export const getPortfolio = async (id)=>{
    return getList(`/portfolios/${id}`)
}

const getList = async (url, params) => {
    return await API_SERVICE.get(url, {
        params: params,
        paramsSerializer: params => {
            return qs.stringify(params, {arrayFormat: "repeat"})
        }
    }).then(response => response.data)
};