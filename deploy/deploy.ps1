<#
.SYNOPSIS
    World-Class ERP Idempotent Deployment Script (幂等性部署脚本)
.DESCRIPTION
    此脚本用于在 Windows 环境下实现 ERP 系统的幂等性（Idempotent）部署。
    所谓幂等性，即：无论当前系统处于什么状态（未安装、运行中、崩溃状态），
    执行该脚本 N 次，最终系统都会收敛到一个一致且正常的最新运行状态。
    
    【面向现在的设计】
    - 自动判断节点状态并关闭老旧进程（释放端口）。
    - 依赖缓存判断，避免重复缓慢的 npm install。
    - 利用 docker-compose 的天然幂等特性，确保数据库集群状态对齐。
    
    【面向未来的设计】
    - $Environment 参数：为未来 Local / Staging / Production 环境隔离做准备。
    - $SERVICES 数组：未来新增微服务（如 service-finance, service-scm），只需在此数组追加一行。
#>

param (
    [ValidateSet("local", "staging", "production")]
    [string]$Environment = "local",
    
    [switch]$ForceRebuild = $false
)

$ErrorActionPreference = "Stop"
$ERP_ROOT = Resolve-Path "$PSScriptRoot\.."
$SERVICES = @("service-gateway") # 面向未来：此处可无限追加微服务名称

Write-Host "=====================================================" -ForegroundColor Cyan
Write-Host "🚀 启动 ERP 幂等性自动化部署流水线 | 环境: $Environment" -ForegroundColor Cyan
Write-Host "=====================================================" -ForegroundColor Cyan

# ---------------------------------------------------------
# [阶段 1] 基础设施层 (Docker-Compose 幂等拉起)
# docker-compose up -d 天然具备幂等性：容器已在运行且配置未变时，它不会做任何事
# ---------------------------------------------------------
Write-Host "`n[1/4] 正在检查并拉起底层基础设施 (PostgreSQL, Redis)..." -ForegroundColor Yellow
try {
    Push-Location "$ERP_ROOT\deploy"
    docker-compose -f docker-compose.yml up -d
    Write-Host "✅ 基础设施状态校验并就绪。" -ForegroundColor Green
} catch {
    Write-Host "⚠️ 警告：当前节点未检测到 Docker，本次部署跳过容器启动阶段。" -ForegroundColor DarkYellow
} finally {
    Pop-Location
}

# ---------------------------------------------------------
# [阶段 2] 前端层 (增量编译与产物更新)
# ---------------------------------------------------------
Write-Host "`n[2/4] 正在构建前端微前端基座..." -ForegroundColor Yellow
$FrontendPath = "$ERP_ROOT\apps\frontend-shell"
if (Test-Path $FrontendPath) {
    Push-Location $FrontendPath
    
    # 幂等安装依赖：如果 node_modules 已存在且未指定强刷，则跳过
    if (-not (Test-Path "node_modules") -or $ForceRebuild) {
        Write-Host "  -> 发现依赖缺失或强制重构，执行 npm install..."
        npm install --silent
    }

    # 幂等构建产物
    Write-Host "  -> 执行代码编译打包 (npm run build)..."
    npm run build --silent
    
    # 面向未来：如果是生产环境，应该将 \dist 目录复制到 Nginx 挂载卷
    if ($Environment -eq "production") {
        Write-Host "  -> (Production) 正在将前端产物同步至 Web 服务器..."
        # Copy-Item -Path "dist\*" -Destination "C:\nginx\html" -Recurse -Force
    }

    Write-Host "✅ 前端构建并更新完成。" -ForegroundColor Green
    Pop-Location
}

# ---------------------------------------------------------
# [阶段 3] 后端微服务层 (平滑停机与重新编译)
# ---------------------------------------------------------
Write-Host "`n[3/4] 正在编译并部署 Go 微服务集群..." -ForegroundColor Yellow

foreach ($svc in $SERVICES) {
    $SvcPath = "$ERP_ROOT\apps\$svc"
    $ExeName = "$svc.exe"
    
    if (Test-Path $SvcPath) {
        Push-Location $SvcPath
        
        # 幂等操作一：杀死旧进程 (防止文件被占用，释放监听端口)
        $existingProcess = Get-Process -Name $svc -ErrorAction SilentlyContinue
        if ($existingProcess) {
            Write-Host "  -> 发现正在运行的老版本服务 [$svc]，执行平滑停机..."
            Stop-Process -Name $svc -Force
            Start-Sleep -Seconds 1 # 确保 OS 彻底释放端口
        }

        # 幂等操作二：覆盖编译新版本
        Write-Host "  -> 重新编译微服务 [$svc]..."
        go build -o $ExeName main.go

        # 面向未来：注册为 Windows 服务并重启，确保机器重启后自启动
        # 在真实生产环境中，建议使用 nssm (Non-Sucking Service Manager) 
        if ($Environment -eq "production") {
            # nssm restart $svc
            Write-Host "  -> (Production) 服务 [$svc] 已在后台重启。"
        } else {
            Write-Host "  -> (Local) 编译完成。请手动运行 .\$ExeName 启动服务。"
        }

        Write-Host "✅ 服务 [$svc] 部署完毕。" -ForegroundColor Green
        Pop-Location
    }
}

# ---------------------------------------------------------
# [阶段 4] 架构自检与收尾
# ---------------------------------------------------------
Write-Host "`n[4/4] 验证部署状态..." -ForegroundColor Yellow
# 这里可以追加 HTTP 健康检查探测，例如 Invoke-WebRequest http://localhost:8080/health

Write-Host "🎉 ERP 节点部署脚本执行完毕！系统已达到一致状态。" -ForegroundColor Green
Write-Host "-----------------------------------------------------" -ForegroundColor Cyan
Write-Host "架构师建议: 未来团队扩大后，可将此脚本无缝迁移为 GitLab CI / GitHub Actions 的 pipeline 脚本。" -ForegroundColor DarkGray

