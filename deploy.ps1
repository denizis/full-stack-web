$ServerIP = "188.132.234.153"
$User = "root"
$RemotePath = "/root/ssh-app"

Write-Host "🚀 Starting Deployment to $ServerIP..." -ForegroundColor Cyan

# 1. Create Remote Directory
Write-Host "📂 Creating remote directory..."
ssh $User@$ServerIP "mkdir -p $RemotePath/frontend"

# 2. Upload Backend Binary
Write-Host "⬆️ Uploading server binary..."
scp .\server-linux $User@$ServerIP:$RemotePath/server

# 3. Upload Frontend
Write-Host "⬆️ Uploading frontend..."
scp -r .\frontend\dist $User@$ServerIP:$RemotePath/frontend/dist

# 4. Start Application
Write-Host "▶️ Starting application on remote server..."
$StartCommand = "cd $RemotePath && chmod +x server && pkill server || true && nohup ./server > server.log 2>&1 &"
ssh $User@$ServerIP $StartCommand

Write-Host "✅ Deployment Complete!" -ForegroundColor Green
Write-Host "🌍 App should be live at: http://$ServerIP:3000"
