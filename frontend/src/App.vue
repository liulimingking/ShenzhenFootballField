<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'

const API_BASE = 'http://localhost:8080'

const fields = ref([])
const districts = ref([])
const loading = ref(true)
const selectedDistrict = ref('')
const selectedStatus = ref('')
const selectedField = ref(null)
const mapContainer = ref(null)
let mapInstance = null
let markers = []

const statusMap = {
  available: { label: '空闲', color: '#4CAF50' },
  busy: { label: '忙碌', color: '#FF9800' },
  full: { label: '已满', color: '#F44336' }
}

const statusOptions = [
  { value: '', label: '全部状态' },
  { value: 'available', label: '空闲' },
  { value: 'busy', label: '忙碌' },
  { value: 'full', label: '已满' }
]

const filteredFields = computed(() => {
  let list = fields.value
  if (selectedDistrict.value) {
    list = list.filter(f => f.district === selectedDistrict.value)
  }
  if (selectedStatus.value) {
    list = list.filter(f => f.status === selectedStatus.value)
  }
  return list
})

async function fetchDistricts() {
  try {
    const res = await axios.get(`${API_BASE}/api/districts`)
    districts.value = ['', ...res.data.districts]
  } catch (e) {
    districts.value = ['', '福田', '南山', '罗湖', '宝安', '龙岗', '龙华', '盐田', '光明', '坪山', '大鹏']
  }
}

async function fetchFields() {
  try {
    const res = await axios.get(`${API_BASE}/api/fields`)
    fields.value = res.data.fields
  } catch (e) {
    loading.value = false
  }
}

async function fetchStatuses() {
  // statuses are static
}

function initMap() {
  if (!mapContainer.value) return
  if (typeof AMap === 'undefined') {
    setTimeout(initMap, 500)
    return
  }
  mapInstance = new AMap.Map(mapContainer.value, {
    zoom: 12,
    center: [114.0630, 22.5480],
    viewMode: '2D',
    mapStyle: 'amap://styles/normal'
  })
}

function renderMarkers() {
  if (!mapInstance) return
  markers.forEach(m => m.setMap(null))
  markers = []

  filteredFields.value.forEach(field => {
    const info = statusMap[field.status] || statusMap.available
    const marker = new AMap.Marker({
      position: new AMap.LngLat(field.lng, field.lat),
      title: field.name,
      content: `<div style="
        background:${info.color};
        width:28px;height:28px;
        border-radius:50%;
        border:2px solid #fff;
        box-shadow:0 2px 6px rgba(0,0,0,0.3);
        display:flex;align-items:center;justify-content:center;
        cursor:pointer;
      "><span style="color:#fff;font-size:14px;font-weight:bold;">⚽</span></div>`,
      offset: new AMap.Pixel(-14, -14)
    })

    marker.on('click', () => {
      selectedField.value = field
      mapInstance.setCenter([field.lng, field.lat])
      mapInstance.setZoom(15)
    })

    marker.setMap(mapInstance)
    markers.push(marker)
  })
}

function selectField(field) {
  selectedField.value = field
  if (mapInstance) {
    mapInstance.setCenter([field.lng, field.lat])
    mapInstance.setZoom(15)
  }
}

function closeDetail() {
  selectedField.value = null
}

function getStatusClass(status) {
  return {
    'status-badge': true,
    'status-available': status === 'available',
    'status-busy': status === 'busy',
    'status-full': status === 'full'
  }
}

onMounted(async () => {
  await Promise.all([fetchDistricts(), fetchFields()])
  loading.value = false

  setTimeout(() => {
    initMap()
    setTimeout(renderMarkers, 1000)
  }, 100)
})

function applyFilters() {
  renderMarkers()
}
</script>

<template>
  <div class="app">
    <header class="header">
      <div class="header-left">
        <span class="logo">⚽</span>
        <h1 class="title">深圳足球场</h1>
      </div>
      <div class="filters">
        <select v-model="selectedDistrict" @change="applyFilters" class="filter-select">
          <option v-for="d in districts" :key="d" :value="d">
            {{ d || '全部区域' }}
          </option>
        </select>
        <select v-model="selectedStatus" @change="applyFilters" class="filter-select">
          <option v-for="s in statusOptions" :key="s.value" :value="s.value">
            {{ s.label }}
          </option>
        </select>
      </div>
    </header>

    <div class="main-content">
      <aside class="sidebar">
        <div class="list-header">
          <span>找到 {{ filteredFields.length }} 个球场</span>
        </div>
        <div v-if="loading" class="loading">加载中...</div>
        <div v-else class="field-list">
          <div
            v-for="field in filteredFields"
            :key="field.id"
            class="field-card"
            :class="{ active: selectedField?.id === field.id }"
            @click="selectField(field)"
          >
            <div class="field-top">
              <h3 class="field-name">{{ field.name }}</h3>
              <span :class="getStatusClass(field.status)">{{ statusMap[field.status]?.label }}</span>
            </div>
            <p class="field-address">{{ field.address }}</p>
            <div class="field-meta">
              <span class="district-tag">{{ field.district }}</span>
              <span class="price">{{ field.priceRange }}</span>
            </div>
          </div>
        </div>
      </aside>

      <div class="map-area">
        <div ref="mapContainer" class="map-container"></div>
        <div v-if="!loading" class="map-hint">地图加载中，请确保已配置高德地图 API Key</div>
      </div>
    </div>

    <Transition name="slide">
      <div v-if="selectedField" class="detail-panel">
        <button class="close-btn" @click="closeDetail">✕</button>
        <div class="detail-header">
          <h2>{{ selectedField.name }}</h2>
          <span :class="getStatusClass(selectedField.status)">
            {{ statusMap[selectedField.status]?.label }}
          </span>
        </div>
        <div class="detail-body">
          <div class="detail-row">
            <span class="detail-label">地址</span>
            <span class="detail-value">{{ selectedField.address }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">区域</span>
            <span class="detail-value">{{ selectedField.district }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">价格</span>
            <span class="detail-value price-highlight">{{ selectedField.priceRange }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">营业时间</span>
            <span class="detail-value">{{ selectedField.openHours }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">联系电话</span>
            <span class="detail-value">
              <a :href="'tel:' + selectedField.phone">{{ selectedField.phone }}</a>
            </span>
          </div>
          <div class="detail-row">
            <span class="detail-label">设施</span>
            <div class="facilities-list">
              <span
                v-for="fac in selectedField.facilities"
                :key="fac"
                class="facility-tag"
              >{{ fac }}</span>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Noto Sans SC', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  background: #F5F5F5;
}

.app {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}

.header {
  background: linear-gradient(135deg, #1B5E20 0%, #2E7D32 100%);
  padding: 0 24px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  z-index: 10;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo {
  font-size: 28px;
}

.title {
  color: #fff;
  font-size: 22px;
  font-weight: 700;
  letter-spacing: 2px;
}

.filters {
  display: flex;
  gap: 10px;
}

.filter-select {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.9);
  color: #333;
  font-size: 14px;
  font-family: inherit;
  cursor: pointer;
  outline: none;
  transition: background 0.2s;
}

.filter-select:hover {
  background: #fff;
}

.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.sidebar {
  width: 360px;
  flex-shrink: 0;
  background: #fff;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #E0E0E0;
  overflow: hidden;
}

.list-header {
  padding: 12px 16px;
  font-size: 13px;
  color: #757575;
  border-bottom: 1px solid #F0F0F0;
  background: #FAFAFA;
}

.loading {
  padding: 32px 16px;
  text-align: center;
  color: #999;
}

.field-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.field-list::-webkit-scrollbar {
  width: 4px;
}

.field-list::-webkit-scrollbar-track {
  background: transparent;
}

.field-list::-webkit-scrollbar-thumb {
  background: #C8C8C8;
  border-radius: 2px;
}

.field-card {
  padding: 14px;
  margin-bottom: 8px;
  background: #FAFAFA;
  border-radius: 10px;
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s;
}

.field-card:hover {
  background: #F0F7F0;
  border-color: #A5D6A7;
}

.field-card.active {
  background: #E8F5E9;
  border-color: #4CAF50;
}

.field-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.field-name {
  font-size: 15px;
  font-weight: 600;
  color: #1B5E20;
  flex: 1;
}

.field-address {
  font-size: 12px;
  color: #757575;
  margin-bottom: 8px;
  line-height: 1.4;
}

.field-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.district-tag {
  font-size: 11px;
  background: #E8F5E9;
  color: #2E7D32;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 500;
}

.price {
  font-size: 12px;
  color: #388E3C;
  font-weight: 600;
}

.status-badge {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 500;
  white-space: nowrap;
  margin-left: 8px;
}

.status-available {
  background: #E8F5E9;
  color: #2E7D32;
}

.status-busy {
  background: #FFF3E0;
  color: #E65100;
}

.status-full {
  background: #FFEBEE;
  color: #C62828;
}

.map-area {
  flex: 1;
  position: relative;
  overflow: hidden;
}

.map-container {
  width: 100%;
  height: 100%;
}

.map-hint {
  position: absolute;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 12px;
  pointer-events: none;
  z-index: 5;
}

.detail-panel {
  position: fixed;
  right: 0;
  top: 60px;
  bottom: 0;
  width: 340px;
  background: #fff;
  box-shadow: -4px 0 16px rgba(0, 0, 0, 0.1);
  z-index: 20;
  overflow-y: auto;
  padding: 24px;
  border-left: 1px solid #E8E8E8;
}

.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 28px;
  height: 28px;
  border: none;
  background: #F5F5F5;
  border-radius: 50%;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}

.close-btn:hover {
  background: #E0E0E0;
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
  padding-right: 40px;
}

.detail-header h2 {
  font-size: 18px;
  color: #1B5E20;
  font-weight: 700;
}

.detail-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.detail-row {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-label {
  font-size: 12px;
  color: #999;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.detail-value {
  font-size: 14px;
  color: #333;
  line-height: 1.5;
}

.detail-value a {
  color: #1B5E20;
  text-decoration: none;
}

.detail-value a:hover {
  text-decoration: underline;
}

.price-highlight {
  font-weight: 600;
  color: #2E7D32 !important;
  font-size: 16px !important;
}

.facilities-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 2px;
}

.facility-tag {
  font-size: 12px;
  padding: 3px 10px;
  background: #E8F5E9;
  color: #2E7D32;
  border-radius: 12px;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateX(100%);
}
</style>
