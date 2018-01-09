#! /bin/bash

echo "Mine being set up"

echo mcb1234 | sudo -S nvidia-smi -pl 105
nvidia-settings -c :0 -a 'GPUFanControlState=1'
nvidia-settings -c :0 -a 'GPUTargetFanSpeed=50'
# nvidia-settings -c :0 -a 'GPUGraphicsClockOffset[3]==100'
nvidia-settings -c :0 -a 'GPUMemoryTransferRateOffset[3]=1150'
