## 用户系统
### 机构信息API
- api: /portal/organization
- method: GET, POST, PUT, DELETE
  - GET: 查询机构信息 返回树形结构信息，支持查询 name
  - PUT: 只允许修改机构名称
### 职位信息API
- api: /portal/position
- method: GET, POST, PUT, DELETE
  - POST: 只允许在主机构下设置职位信息，子机构不允许设置职位信息
  - PUT: 只允许修改职位名称