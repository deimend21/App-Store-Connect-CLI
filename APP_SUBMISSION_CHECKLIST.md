# iOS App Store 上架提审完全指南与标准清单 (SOP)

> 本文档整理自实际 App 上架提审实战经验，旨在帮助后续所有 iOS 项目以最高的效率、零踩坑完成从**项目配置、证书打包、元数据填写、素材上传到最终提交审核**的全流程。

---

## 目录
- [一、提审必备准备清单（提前备齐）](#一提审必备准备清单提前备齐)
- [二、五大核心避坑要点（极其重要）](#二五大核心避坑要点极其重要)
- [三、极速流水线命令速查（按顺序执行）](#三极速流水线命令速查按顺序执行)
- [四、高频报错与一键解决方案](#四高频报错与一键解决方案)

---

## 一、提审必备准备清单（提前备齐）

在开始使用 `asc` 命令行前，建议将以下内容整理为一个清单文件或记事本：

### 1. 基础产品与标识信息
- **Bundle ID**：例如 `com.yourcompany.appname`（唯一且全小写/点分）。
- **App 官方名称（Store Name）**：最多 30 个字符，例如 `钓鱼相机 - 对照拍同款 (DuoCam)`。
- **副标题（Subtitle）**：最多 30 个字符，例如 `洋葱皮半透明对照，轻松拍同款`。
- **桌面显示名称（CFBundleDisplayName）**：手机图标下方文字，通常 4~6 个汉字，例如 `钓鱼相机`。
- **SKU**：后台内部唯一编号，例如 `duocam-pose-match`。
- **分类**：主分类（如 `PHOTO_AND_VIDEO`）、副分类（如 `UTILITIES`）。
- **定价与分发**：免费（Free, Tier 0）或指定定价等级，覆盖国家/地区（默认全部）。

### 2. 两个公开合法网址（必须公开可用）
- **隐私政策网址（Privacy Policy URL）**：必须是公开可访问的 HTTP/HTTPS 页面。
- **技术支持 / 营销网址（Support / Marketing URL）**：必须公开可访问，包含联系方式或使用说明。
> 💡 **推荐托管方式**：
> - 公开的 Notion 页面（发布为网页并复制 Web 链接，如 `https://xxx.notion.site/...`）
> - GitHub Pages（使用公开仓库开通 Pages 功能）
> - 个人或公司独立域名网页
> 
> ⚠️ **严禁**：私有 GitHub 仓库（`github.com/user/private-repo` 会返回 404）、需登录才能访问的飞书/腾讯文档。

### 3. 审核人联系信息（Review Details）
- **姓名**：联系人拼音或英文（如 `Jiahua Lin`）。
- **邮箱**：接收审核反馈的邮箱（如 `your_email@domain.com`）。
- **电话号码**：**必须包含国际区号**，例如 `+86 152xxxxxxxx`（否则报参数格式错误）。
- **Demo 账号与密码**：若 App 包含登录体系，必须提供有效的测试账号；若无登录体系，注明免登录。
- **审核备注（Notes）**：若有特殊硬件依赖、蓝牙、相机对照等功能，简要说明测试方式。

### 4. 宣传截图素材（Screenshots）
苹果强制要求至少上传对应平台主流机型的截图（推荐 3 ~ 5 张）：
- **iPhone（6.7 英寸超视网膜）**：分辨率严格为 **`1290 × 2796`**（竖屏）或 **`2796 × 1290`**（横屏）。
- **iPad（12.9 英寸 iPad Pro）**：如果项目支持 iPad（Universal），必须提供 **`2048 × 2732`**（竖屏）或 **`2732 × 2048`**（横屏）。
> 格式要求：PNG 或高画质 JPG，RGB 色彩空间，无透明通道（Alpha）。

### 5. 工程配置与权限文案（Xcode / project.yml）
- **权限描述（Usage Descriptions）**：如果用到了相机、相册、麦克风等，必须在 `Info.plist` 提供详尽的人性化说明：
  - `NSCameraUsageDescription`
  - `NSPhotoLibraryUsageDescription` / `NSPhotoLibraryAddUsageDescription`
  - `NSMicrophoneUsageDescription`
- **全屏启动屏声明**：
  - Xcode 项目中必须配置有效的 Launch Screen（或在 `project.yml` 中加上 `UILaunchScreen: {}`），并且启用 `UIRequiresFullScreen: true`，否则 iPad 审核会被拦截。

---

## 二、五大核心避坑要点（极其重要）

### 坑 1：App 数据隐私问卷（App Privacy）必须在网页端签署发布
- **现象**：执行 `asc review submit` 时报错：
  `Associated errors for /v1/appDataUsages/: You must have published answers to your app's data usages.`
- **原因**：苹果官方规定，涉及法律免责的“数据隐私问卷”无法通过公开 API Key 签署，必须由开发者主账号在网页完成。
- **秒解法**：
  1. 打开 `https://appstoreconnect.apple.com/apps/<APP_ID>/appPrivacy`
  2. 点击“开始”，选择 **“否，我们不会从此 App 中收集数据”**（单机或不追踪 App）。
  3. 点击“存储”，最后点击右上角蓝色的 **“发布”** 按钮即可。

### 坑 2：隐私政策与技术支持网址不可用（404 必拒）
- 苹果审核团队会真机点击你的 Support URL 和 Privacy URL。
- 私有仓库链接（返回 404）或需要登录账号的页面会直接被判定违反 **Guideline 5.1.1 (Data Collection and Storage)**。
- 提审前先在无痕浏览器中测试链接是否能直接打开。

### 坑 3：审核联系电话缺少国际区号
- 填写 Review Details 时，中国大陆手机号如果只填 `152xxxxxxxx`，API 会报错或被苹果后台拒绝。
- **必须填写标准国际格式**：`+86 152xxxxxxxx`。

### 坑 4：IPA 架构与描述文件（Provisioning Profile）类型错误
- 上传 App Store 的 IPA 必须使用 **Apple Distribution 证书** 和 **App Store 类型的 Provisioning Profile** 签名。
- 使用 Development 或 Ad-Hoc 证书构建的包会在上传或校验时报 `ITMS-90034 / ITMS-90161` 错误。

### 坑 5：IPA 上传后需等待后台处理（VALID）
- `asc builds upload` 上传成功后，苹果后台通常需要 **5 ~ 15 分钟** 进行静态安全扫描。
- 此时 build 处于 `PROCESSING` 状态，强行提交会导致无法关联。
- 使用 `asc builds wait --build-id <BUILD_ID>` 或 `asc builds info` 确认状态变为 `VALID` 后再提审。

---

## 三、极速流水线命令速查（按顺序执行）

以下是一套经过验证的端到端 CLI 执行模板：

### 步骤 1：检查账号与环境
```bash
# 检查默认登录 Profile 或验证有效性
asc auth status
```

### 步骤 2：创建 Bundle ID 与描述文件（仅首发需要）
```bash
# 注册 Bundle ID
asc bundle-ids create --identifier "com.example.myapp" --name "MyApp" --platform IOS

# 创建 App Store 发布描述文件
asc profiles create \
  --name "MyApp AppStore Profile" \
  --bundle-id "com.example.myapp" \
  --profile-type IOS_APP_STORE \
  --certificate-id "<DISTRIBUTION_CERT_ID>" \
  --output-path ~/Library/MobileDevice/Provisioning\ Profiles/
```

### 步骤 3：创建 App 与版本（仅首发需要）
```bash
# 创建 App
asc apps create \
  --name "App完整标题 (最多30字)" \
  --bundle-id "com.example.myapp" \
  --sku "myapp-sku" \
  --locale "zh-Hans" \
  --primary-category "PHOTO_AND_VIDEO"

# 创建 1.0 版本
asc versions create \
  --app "<APP_ID>" \
  --version "1.0" \
  --platform IOS
```

### 步骤 4：配置元数据与公开网址
```bash
# 1. 配置 App 级隐私政策与副标题 (app-info)
asc localizations update \
  --type app-info \
  --id "<APP_INFO_LOCALIZATION_ID>" \
  --subtitle "副标题说明" \
  --privacy-policy-url "https://flying-foe-91b.notion.site/your-privacy-page"

# 2. 配置版本级描述、关键词与技术支持网址
asc localizations update \
  --id "<VERSION_LOCALIZATION_ID>" \
  --description "详细的功能介绍、使用方法..." \
  --keywords "关键词1,关键词2,关键词3" \
  --support-url "https://flying-foe-91b.notion.site/your-support-page"
```

### 步骤 5：配置审核联系人信息
```bash
asc review-details update \
  --id "<REVIEW_DETAILS_ID>" \
  --first-name "Jiahua" \
  --last-name "Lin" \
  --phone-number "+86 15217870075" \
  --email "your_email@domain.com" \
  --notes "如：本应用为对照摄影工具，打开即可直接体验，无需登录账号。"
```

### 步骤 6：上传宣传截图
```bash
# 1. 上传 iPhone 6.7 寸截图（1290x2796）
asc screenshots upload \
  --version-localization-id "<VERSION_LOCALIZATION_ID>" \
  --display-type APP_IPHONE_67 \
  --path "./screenshots/iphone"

# 2. 上传 iPad Pro 12.9 寸截图（2048x2732，通用应用必须）
asc screenshots upload \
  --version-localization-id "<VERSION_LOCALIZATION_ID>" \
  --display-type APP_IPAD_PRO_3GEN_129 \
  --path "./screenshots/ipad"
```

### 步骤 7：上传 IPA 并等待解析完成
```bash
# 上传构建包
asc builds upload \
  --app "<APP_ID>" \
  --ipa "./build/MyApp.ipa" \
  --platform IOS

# 查看最近的构建状态（等待 processingState 变成 VALID）
asc builds list --app "<APP_ID>" --sort "-uploadedDate" --limit 1
```

### 步骤 8：网页端发布 App Privacy（参考避坑要点 1）
前往：`https://appstoreconnect.apple.com/apps/<APP_ID>/appPrivacy` 完成发布。

### 步骤 9：正式提交审核并验证状态
```bash
# 提交审核
asc review submit \
  --app "<APP_ID>" \
  --version "1.0" \
  --build-id "<BUILD_ID>" \
  --confirm

# 验证状态（看到 WAITING_FOR_REVIEW 且 blockerCount: 0 即大功告成）
asc status --app "<APP_ID>"
```

---

## 四、高频报错与一键解决方案

| 报错信息 / 现象 | 产生原因 | 一键解决办法 |
| :--- | :--- | :--- |
| `You must have published answers to your app's data usages` | 苹果要求开发者账号签署 App Privacy 问卷 | 打开 `https://appstoreconnect.apple.com/apps/<APP_ID>/appPrivacy`，选“否，不收集数据”并点击右上角“发布” |
| `phone number is invalid` | 审核联系人电话未带国际区号 | 改为 `+86 1xxxxxxxxxx` 格式 |
| `bundleId is already reserved or exists` | Bundle ID 已被占用或已被其他团队注册 | 更换 Bundle ID 字符串重新注册 |
| `Screenshot dimension mismatch` | 截图尺寸不符合苹果固定规格 | 严格导出为 `1290x2796` (iPhone) 或 `2048x2732` (iPad) |
| `Version is not in valid state / Cannot be reviewed` | 缺少必要的元数据、截图或构建包处理中 | 运行 `asc validate --app <APP_ID> --version <VERSION>` 打印缺失项 |
| `Review submission has items with errors` | 旧的审核提交残留了冲突的 Draft 记录 | 使用 `asc review cancel` 或在网页后台丢弃旧草稿再提审 |
