const api = require('../services/api')

// 配置缓存
let configCache = {}
let cacheTimestamp = 0
const CACHE_EXPIRE = 3600000 // 1小时缓存

/**
 * 获取配置值（带缓存）
 * @param {string} name - 配置名称
 * @param {boolean} useCache - 是否使用缓存
 * @returns {Promise<string>} 配置值
 */
async function getConfig(name, useCache = true) {
  // 尝试从缓存读取
  if (useCache && configCache[name] && (Date.now() - cacheTimestamp < CACHE_EXPIRE)) {
    return configCache[name]
  }

  try {
    const res = await api.getConfig(name)
    if (res.code === 0 && res.data) {
      const value = res.data.value
      
      // 更新缓存
      configCache[name] = value
      cacheTimestamp = Date.now()
      
      return value
    }
  } catch (err) {
    console.error(`获取配置 ${name} 失败`, err)
  }
  
  return null
}

/**
 * 批量获取配置（带缓存）
 * @param {string[]} names - 配置名称数组
 * @param {boolean} useCache - 是否使用缓存
 * @returns {Promise<Object>} 配置对象 { name: value }
 */
async function getConfigs(names, useCache = true) {
  // 检查缓存
  if (useCache) {
    const cached = {}
    let allCached = true
    
    for (const name of names) {
      if (configCache[name] && (Date.now() - cacheTimestamp < CACHE_EXPIRE)) {
        cached[name] = configCache[name]
      } else {
        allCached = false
        break
      }
    }
    
    if (allCached) {
      return cached
    }
  }

  try {
    const res = await api.getConfigsByNames(names)
    if (res.code === 0 && res.data) {
      // 更新缓存
      Object.assign(configCache, res.data)
      cacheTimestamp = Date.now()
      
      return res.data
    }
  } catch (err) {
    console.error('批量获取配置失败', err)
  }
  
  return {}
}

/**
 * 设置导航栏标题
 * @param {string} configName - 配置名称
 * @param {string} defaultTitle - 默认标题（配置获取失败时使用）
 */
async function setNavigationBarTitle(configName, defaultTitle = '') {
  try {
    const value = await getConfig(configName)
    if (value) {
      wx.setNavigationBarTitle({
        title: value
      })
    } else if (defaultTitle) {
      wx.setNavigationBarTitle({
        title: defaultTitle
      })
    }
  } catch (err) {
    console.error('设置导航栏标题失败', err)
    if (defaultTitle) {
      wx.setNavigationBarTitle({
        title: defaultTitle
      })
    }
  }
}

/**
 * 清除配置缓存
 */
function clearConfigCache() {
  configCache = {}
  cacheTimestamp = 0
}

module.exports = {
  getConfig,
  getConfigs,
  setNavigationBarTitle,
  clearConfigCache
}
