const api = require('../../utils/request.js')

Page({
  data: {
    id: null,
    form: {
      consignee: '',
      phone: '',
      province: '',
      city: '',
      district: '',
      address: '',
      postal_code: '',
      is_default: 0
    },
    regions: ['北京市', '天津市', '河北省', '山西省', '内蒙古', '辽宁省', '吉林省', '黑龙江省', '上海市', '江苏省', '浙江省', '安徽省', '福建省', '江西省', '山东省', '河南省', '湖北省', '湖南省', '广东省', '广西省', '海南省', '重庆市', '四川省', '贵州省', '云南省', '西藏', '陕西省', '甘肃省', '青海省', '宁夏', '新疆', '台湾省', '香港', '澳门'],
    selectedRegion: ''
  },

  onLoad(options) {
    if (options.id) {
      this.setData({ id: parseInt(options.id) })
      this.loadAddress()
    }
  },

  async loadAddress() {
    try {
      const res = await api.getAddress(this.data.id)
      const address = res.data
      this.setData({
        form: address,
        selectedRegion: address.province
      })
    } catch (err) {
      console.error(err)
    }
  },

  onInput(e) {
    const { field } = e.currentTarget.dataset
    this.setData({
      [`form.${field}`]: e.detail.value
    })
  },

  onRegionChange(e) {
    const { value } = e.detail
    this.setData({
      selectedRegion: value,
      'form.province': value,
      'form.city': value,
      'form.district': value
    })
  },

  onDefaultChange(e) {
    this.setData({
      'form.is_default': e.detail.value ? 1 : 0
    })
  },

  async onSubmit() {
    const { form } = this.data

    if (!form.consignee) {
      wx.showToast({ title: '请输入收货人', icon: 'none' })
      return
    }
    if (!form.phone) {
      wx.showToast({ title: '请输入手机号', icon: 'none' })
      return
    }
    if (!form.province) {
      wx.showToast({ title: '请选择地区', icon: 'none' })
      return
    }
    if (!form.address) {
      wx.showToast({ title: '请输入详细地址', icon: 'none' })
      return
    }

    try {
      if (this.data.id) {
        await api.updateAddress(this.data.id, form)
      } else {
        await api.createAddress(form)
      }
      wx.showToast({ title: '保存成功', icon: 'success' })
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)
    } catch (err) {
      console.error(err)
    }
  }
})
