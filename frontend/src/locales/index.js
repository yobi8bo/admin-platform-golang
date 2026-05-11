import { createI18n } from 'vue-i18n'

const messages = {
  'zh-CN': {
    app: {
      name: 'Admin Platform',
      console: '企业控制台',
      dashboard: '仪表盘',
      profile: '个人中心',
      logout: '退出登录'
    },
    common: {
      actions: '操作',
      createdAt: '时间',
      delete: '删除',
      download: '下载',
      edit: '编辑',
      id: 'ID',
      message: '消息',
      reset: '重置',
      save: '保存',
      search: '查询',
      status: '状态',
      success: '成功',
      failed: '失败'
    },
    menu: {
      Dashboard: '仪表盘',
      System: '系统管理',
      User: '用户管理',
      Role: '角色管理',
      Menu: '菜单管理',
      Dept: '部门管理',
      File: '文件管理',
      Audit: '审计日志',
      LoginLog: '登录日志',
      OperationLog: '操作日志'
    },
    audit: {
      loginLog: '登录日志',
      loginLogDesc: '查看账号登录结果、来源 IP 与提示信息',
      operationLog: '操作日志',
      operationLogDesc: '记录用户在系统内触发的关键 API 操作',
      username: '用户名',
      user: '用户',
      operation: '操作',
      statusCode: '状态码'
    },
    login: {
      heading: '企业后台管理平台',
      tagline: '权限、文件、审计与组织数据集中管理。',
      title: '登录',
      desc: '使用系统账号进入管理控制台',
      username: '用户名',
      password: '密码',
      usernameRequired: '请输入用户名',
      passwordRequired: '请输入密码',
      submit: '登录',
      demo: '默认账号：admin / Admin@123',
      failed: '登录失败'
    },
    dashboard: {
      title: '仪表盘',
      desc: '系统运行入口与关键模块概览',
      refresh: '刷新权限数据',
      menuCount: '菜单项',
      permissionCount: '按钮权限',
      roleCount: '当前角色',
      currentUser: '登录用户',
      authReady: '认证与个人资料接口已接入',
      systemReady: '系统管理模块使用后端权限控制',
      fileReady: '文件上传、下载链接与删除操作已对齐',
      auditReady: '审计日志按分页表格展示'
    },
    file: {
      title: '文件管理',
      desc: '上传文件，生成临时下载链接',
      upload: '上传文件',
      keyword: '文件名 / 类型',
      name: '文件名',
      type: '类型',
      size: '大小',
      uploadSuccess: '上传成功',
      deleteConfirm: '确认删除该文件？',
      deleted: '已删除'
    }
  },
  'en-US': {
    app: {
      name: 'Admin Platform',
      console: 'Enterprise Console',
      dashboard: 'Dashboard',
      profile: 'Profile',
      logout: 'Logout'
    },
    common: {
      actions: 'Actions',
      createdAt: 'Time',
      delete: 'Delete',
      download: 'Download',
      edit: 'Edit',
      id: 'ID',
      message: 'Message',
      reset: 'Reset',
      save: 'Save',
      search: 'Search',
      status: 'Status',
      success: 'Success',
      failed: 'Failed'
    },
    menu: {
      Dashboard: 'Dashboard',
      System: 'System',
      User: 'Users',
      Role: 'Roles',
      Menu: 'Menus',
      Dept: 'Departments',
      File: 'Files',
      Audit: 'Audit',
      LoginLog: 'Login Logs',
      OperationLog: 'Operation Logs'
    },
    audit: {
      loginLog: 'Login Logs',
      loginLogDesc: 'Review account login results, source IPs, and messages',
      operationLog: 'Operation Logs',
      operationLogDesc: 'Track key API operations triggered by users',
      username: 'Username',
      user: 'User',
      operation: 'Operation',
      statusCode: 'Status Code'
    },
    login: {
      heading: 'Enterprise Admin Platform',
      tagline: 'Centralized access control, files, audit trails, and organization data.',
      title: 'Sign in',
      desc: 'Use your system account to enter the console',
      username: 'Username',
      password: 'Password',
      usernameRequired: 'Please enter username',
      passwordRequired: 'Please enter password',
      submit: 'Sign in',
      demo: 'Default account: admin / Admin@123',
      failed: 'Sign in failed'
    },
    dashboard: {
      title: 'Dashboard',
      desc: 'System entry point and module overview',
      refresh: 'Refresh permission data',
      menuCount: 'Menus',
      permissionCount: 'Button Permissions',
      roleCount: 'Roles',
      currentUser: 'Current User',
      authReady: 'Authentication and profile APIs are connected',
      systemReady: 'System management uses backend permission control',
      fileReady: 'File upload, download links, and deletion are aligned',
      auditReady: 'Audit logs are displayed in paginated tables'
    },
    file: {
      title: 'Files',
      desc: 'Upload files and generate temporary download links',
      upload: 'Upload File',
      keyword: 'Filename / type',
      name: 'Filename',
      type: 'Type',
      size: 'Size',
      uploadSuccess: 'Uploaded',
      deleteConfirm: 'Delete this file?',
      deleted: 'Deleted'
    }
  }
}

export default createI18n({
  legacy: false,
  locale: localStorage.getItem('locale') || 'zh-CN',
  fallbackLocale: 'zh-CN',
  messages
})
