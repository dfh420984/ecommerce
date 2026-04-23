import * as XLSX from 'xlsx'

/**
 * 通用Excel导出函数
 * @param {Array} data - 要导出的数据数组
 * @param {String} fileName - 导出的文件名（不含扩展名）
 * @param {String} sheetName - 工作表名称（默认：Sheet1）
 * @param {Object} columnMap - 列映射配置，用于自定义列标题和顺序
 */
export function exportToExcel(data, fileName, sheetName = 'Sheet1', columnMap = null) {
  if (!data || data.length === 0) {
    throw new Error('没有可导出的数据')
  }

  let exportData = data

  // 如果提供了列映射，转换数据格式
  if (columnMap) {
    exportData = data.map(item => {
      const newItem = {}
      Object.keys(columnMap).forEach(key => {
        newItem[columnMap[key]] = item[key] !== undefined ? item[key] : ''
      })
      return newItem
    })
  }

  // 创建工作簿和工作表
  const worksheet = XLSX.utils.json_to_sheet(exportData)
  const workbook = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(workbook, worksheet, sheetName)

  // 生成Excel文件并下载
  XLSX.writeFile(workbook, `${fileName}.xlsx`)
}

/**
 * 从URL参数构建查询字符串
 * @param {Object} params - 参数对象
 * @returns {String} 查询字符串
 */
export function buildQueryString(params) {
  const query = []
  for (const key in params) {
    if (params[key] !== undefined && params[key] !== null && params[key] !== '') {
      query.push(`${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`)
    }
  }
  return query.join('&')
}
