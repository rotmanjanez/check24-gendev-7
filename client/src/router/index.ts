import { createRouter, createWebHashHistory } from 'vue-router'
import ProductSearchView from '@/views/ProductSearchView.vue'
import SharedProductsView from '@/views/SharedProductsView.vue'

const routes = [
    { path: '/', name: 'products', component: ProductSearchView },
    { path: '/internetproducts/share', name: 'shared-products', component: SharedProductsView },
]

export const router = createRouter({
    history: createWebHashHistory(import.meta.env.BASE_URL),
    routes,
})