import { createRouter, createWebHashHistory } from 'vue-router'
import ProductSearchView from '@/views/ProductSearchView.vue'

const routes = [
    { path: '/', name: 'products', component: ProductSearchView },
]

export const router = createRouter({
    history: createWebHashHistory(import.meta.env.BASE_URL),
    routes,
})