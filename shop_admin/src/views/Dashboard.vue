<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card" v-loading="loading">
          <div class="stat-icon" style="background: #409EFF;">
            <el-icon :size="30"><User /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.today_users || 0 }}</div>
            <div class="stat-label">今日新增用户</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" v-loading="loading">
          <div class="stat-icon" style="background: #67C23A;">
            <el-icon :size="30"><List /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.today_orders || 0 }}</div>
            <div class="stat-label">今日订单数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" v-loading="loading">
          <div class="stat-icon" style="background: #E6A23C;">
            <el-icon :size="30"><Money /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">¥{{ formatMoney(stats.today_sales) }}</div>
            <div class="stat-label">今日销售额</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" v-loading="loading">
          <div class="stat-icon" style="background: #F56C6C;">
            <el-icon :size="30"><ShoppingCart /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.pending_orders || 0 }}</div>
            <div class="stat-label">待处理订单</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 销售趋势和热销商品 -->
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>销售趋势（最近7天）</span>
          </template>
          <div ref="salesChartRef" style="height: 300px;"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>热销商品TOP10</span>
          </template>
          <el-table :data="stats.hot_products || []" style="width: 100%" :max-height="300">
            <el-table-column prop="name" label="商品名称" show-overflow-tooltip />
            <el-table-column prop="sales" label="销量" width="80" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 用户增长趋势 -->
    <el-row style="margin-top: 20px;">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>用户增长趋势（最近7天）</span>
          </template>
          <div ref="usersChartRef" style="height: 300px;"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { User, List, Money, ShoppingCart } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import request from '@/utils/request'

const loading = ref(false)
const salesChartRef = ref(null)
const usersChartRef = ref(null)

const stats = reactive({
  today_sales: 0,
  yesterday_sales: 0,
  month_sales: 0,
  today_users: 0,
  today_orders: 0,
  pending_orders: 0,
  hot_products: []
})

// 获取仪表盘统计数据
const fetchDashboardStats = async () => {
  loading.value = true
  try {
    const res = await request.get('/admin/statistics/dashboard')
    Object.assign(stats, res.data)
  } catch (error) {
    console.error('获取统计数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 渲染销售趋势图
const renderSalesChart = async () => {
  try {
    const res = await request.get('/admin/statistics/sales-trend')
    const data = res.data
    
    if (!salesChartRef.value) return
    
    const chart = echarts.init(salesChartRef.value)
    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: { type: 'cross' }
      },
      legend: {
        data: ['销售额', '订单数']
      },
      xAxis: {
        type: 'category',
        data: data.map(item => item.date)
      },
      yAxis: [
        {
          type: 'value',
          name: '销售额(元)',
          position: 'left'
        },
        {
          type: 'value',
          name: '订单数',
          position: 'right'
        }
      ],
      series: [
        {
          name: '销售额',
          type: 'line',
          data: data.map(item => item.sales),
          smooth: true,
          itemStyle: { color: '#409EFF' }
        },
        {
          name: '订单数',
          type: 'bar',
          yAxisIndex: 1,
          data: data.map(item => item.orders),
          itemStyle: { color: '#67C23A' }
        }
      ]
    }
    chart.setOption(option)
    
    // 响应式
    window.addEventListener('resize', () => chart.resize())
  } catch (error) {
    console.error('获取销售趋势失败:', error)
  }
}

// 渲染用户增长趋势图
const renderUsersChart = async () => {
  try {
    const res = await request.get('/admin/statistics/users-trend')
    const data = res.data
    
    if (!usersChartRef.value) return
    
    const chart = echarts.init(usersChartRef.value)
    const option = {
      tooltip: {
        trigger: 'axis'
      },
      xAxis: {
        type: 'category',
        data: data.map(item => item.date)
      },
      yAxis: {
        type: 'value',
        name: '新增用户数'
      },
      series: [
        {
          name: '新增用户',
          type: 'line',
          data: data.map(item => item.count),
          smooth: true,
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(103, 194, 58, 0.3)' },
              { offset: 1, color: 'rgba(103, 194, 58, 0.1)' }
            ])
          },
          itemStyle: { color: '#67C23A' }
        }
      ]
    }
    chart.setOption(option)
    
    // 响应式
    window.addEventListener('resize', () => chart.resize())
  } catch (error) {
    console.error('获取用户趋势失败:', error)
  }
}

// 格式化金额
const formatMoney = (value) => {
  if (!value) return '0.00'
  return Number(value).toFixed(2).replace(/\B(?=(\d{3})+(?!\d))/g, ',')
}

onMounted(async () => {
  await fetchDashboardStats()
  await nextTick()
  renderSalesChart()
  renderUsersChart()
})
</script>

<style scoped>
.dashboard {
  padding: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #999;
  margin-top: 5px;
}
</style>
