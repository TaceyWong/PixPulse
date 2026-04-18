<script lang="ts" setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { Convert, SelectFile, OpenDirectory, GetFileBase64 } from '../wailsjs/go/main/App';
import * as runtime from '../wailsjs/runtime/runtime';
import pulseLogo from './assets/images/pulse.svg';

type Mode = 'color' | 'bw' | 'render';

const inputPath = ref('');
const inputExt = ref('');
const outputPath = ref('');
const inputBase64 = ref('');
const outputBase64 = ref('');
const mode = ref<Mode>('color');
const isProcessing = ref(false);
const status = ref('');
const error = ref('');
const isDragOver = ref(false);

const updatePreviews = async () => {
  if (inputPath.value) {
    try {
      inputBase64.value = await GetFileBase64(inputPath.value);
    } catch (e) {
      console.error("预览加载失败", e);
    }
  }
  if (status.value === '转换成功' && outputPath.value) {
    try {
      outputBase64.value = await GetFileBase64(outputPath.value);
    } catch (e) {
      console.error("输出加载失败", e);
    }
  } else {
    outputBase64.value = '';
  }
};

const processInput = async (input: string) => {
  if (!input) return;
  const lastDot = input.lastIndexOf('.');
  const ext = input.substring(lastDot).toLowerCase();
  
  if (!['.png', '.svg'].includes(ext)) {
    error.value = "仅支持 PNG 和 SVG 格式";
    return;
  }

  inputPath.value = input;
  inputExt.value = ext;
  error.value = '';
  status.value = '';
  outputBase64.value = '';

  if (ext === '.svg') {
    mode.value = 'render';
  } else {
    mode.value = 'color';
  }

  const targetExt = mode.value === 'render' ? '.png' : '.svg';
  outputPath.value = input.substring(0, lastDot) + targetExt;
  await updatePreviews();
};

const handleDrop = (x: number, y: number, paths: string[]) => {
  isDragOver.value = false;
  if (paths && paths.length > 0) {
    processInput(paths[0]);
  }
};

onMounted(() => {
  setTimeout(() => {
    runtime.OnFileDrop(handleDrop, true);
  }, 300);
});

onUnmounted(() => {
  runtime.OnFileDropOff();
});

const handleSelectFile = async () => {
  try {
    const selected = await SelectFile();
    if (selected) await processInput(selected);
  } catch (e: any) {
    error.value = "选择失败: " + e;
  }
};

const startConversion = async () => {
  if (!inputPath.value) return;
  isProcessing.value = true;
  status.value = '执行中...';
  error.value = '';
  outputBase64.value = '';

  try {
    const resp = await Convert({
      inputPath: inputPath.value,
      outputPath: outputPath.value,
      mode: mode.value
    });

    if (resp && resp.success) {
      status.value = '转换成功';
      await updatePreviews();
    } else {
      error.value = resp ? resp.error : '转换失败';
      status.value = '转换失败';
    }
  } catch (e: any) {
    error.value = e.toString();
    status.value = '转换失败';
  } finally {
    isProcessing.value = false;
  }
};

const handleCopyImage = async () => {
  if (!outputBase64.value) return;
  try {
    const response = await fetch(outputBase64.value);
    const blob = await response.blob();
    const isSvg = blob.type === 'image/svg+xml';

    if (isSvg) {
      const svgCode = await blob.text();
      const item = new ClipboardItem({
        'text/plain': new Blob([svgCode], { type: 'text/plain' }),
        'image/svg+xml': blob
      });
      await navigator.clipboard.write([item]);
    } else {
      const item = new ClipboardItem({ [blob.type]: blob });
      await navigator.clipboard.write([item]);
    }
    
    const oldStatus = status.value;
    status.value = isSvg ? '矢量数据已存入剪贴板' : '图片已存入剪贴板';
    setTimeout(() => (status.value = oldStatus), 2000);
  } catch (e: any) {
    error.value = "复制失败: " + e.message;
  }
};

const setMode = (m: Mode) => {
  if (!inputPath.value) return;
  mode.value = m;
  const lastDot = inputPath.value.lastIndexOf('.');
  const targetExt = m === 'render' ? '.png' : '.svg';
  outputPath.value = inputPath.value.substring(0, lastDot) + targetExt;
};

const minimize = () => runtime.WindowMinimise();
const toggleMaximize = () => runtime.WindowToggleMaximise();
const quit = () => runtime.Quit();
</script>

<template>
  <div class="app-wrapper" style="--wails-drop-target: drop">
    <nav class="title-bar" style="--wails-draggable: drag">
      <div class="brand-container">
        <img :src="pulseLogo" class="logo" alt="logo" />
        <div class="brand mono uppercase-tracked">PIXPULSE</div>
      </div>
      
      <div class="window-controls" style="--wails-draggable: none">
        <button class="win-btn maximize" title="全屏/恢复" @click="toggleMaximize">
          <span class="icon"></span>
        </button>
        <button class="win-btn minimize" title="最小化" @click="minimize">
          <span class="icon"></span>
        </button>
        <button class="win-btn close" title="关闭" @click="quit">
          <span class="icon"></span>
        </button>
      </div>
    </nav>

    <main class="container">
      <section class="main-layout">
        <div class="preview-panel" :class="{ 'empty': !inputPath }">
          <div class="panel-header mono small">SOURCE ({{ inputExt.toUpperCase().replace('.','') || 'EMPTY' }})</div>
          <div class="preview-content" @click="handleSelectFile">
            <img v-if="inputBase64" :src="inputBase64" class="preview-img" />
            <div v-else class="placeholder">
              <span class="mono">点击或拖入</span>
              <span class="text-muted tiny">PNG / SVG</span>
            </div>
          </div>
          <div class="panel-actions-placeholder"></div>
        </div>

        <div class="center-controls">
          <div class="mode-stack">
             <template v-if="inputExt === '.png' || !inputPath">
                <button class="mode-btn mono" :class="{ active: mode === 'color' }" :disabled="!inputPath" @click="setMode('color')">彩色</button>
                <button class="mode-btn mono" :class="{ active: mode === 'bw' }" :disabled="!inputPath" @click="setMode('bw')">黑白</button>
             </template>
             <template v-else-if="inputExt === '.svg'">
                <button class="mode-btn mono" :class="{ active: mode === 'render' }" @click="setMode('render')">渲染</button>
             </template>
          </div>
          <button class="execute-btn mono" :disabled="isProcessing || !inputPath" @click="startConversion">
            {{ isProcessing ? '...' : '转换' }}
          </button>
        </div>

        <div class="preview-panel" :class="{ 'empty': !outputBase64 }">
          <div class="panel-header mono small">RESULT</div>
          <div class="preview-content">
            <img v-if="outputBase64" :src="outputBase64" class="preview-img" />
            <div v-else class="placeholder"><span class="mono text-muted">等待执行</span></div>
          </div>
          <div class="panel-actions">
            <template v-if="outputBase64">
              <!-- 修复：根据模式动态显示文字 -->
              <button class="ghost-btn mono small" @click="handleCopyImage">
                {{ mode === 'render' ? '复制图片' : '复制矢量' }}
              </button>
              <button class="ghost-btn mono small" @click="OpenDirectory(outputPath)">目录</button>
            </template>
            <div v-else class="actions-empty"></div>
          </div>
        </div>
      </section>

      <footer class="status-footer">
        <div class="status-wrap" :class="{ 'is-error': error }">
          <span class="mono small">{{ error || status || '就绪' }}</span>
        </div>
        <div v-if="inputPath" class="path-info mono tiny text-muted">{{ inputPath }}</div>
      </footer>
    </main>
  </div>
</template>

<style scoped>
.app-wrapper { display: flex; flex-direction: column; height: 100vh; background-color: var(--bg-color); border: 1px solid var(--border-default); user-select: none; }
.title-bar { height: 40px; display: flex; align-items: center; justify-content: space-between; padding: 0 16px; background: rgba(255, 255, 255, 0.03); }
.brand-container {
  display: flex;
  align-items: center;
  gap: 10px;
}
.logo {
  height: 18px;
  width: auto;
  filter: drop-shadow(0 0 2px rgba(255,255,255,0.3));
}
.brand { font-size: 11px; letter-spacing: 3px; opacity: 0.6; }
.window-controls { display: flex; gap: 8px; align-items: center; }
.win-btn { width: 12px; height: 12px; border-radius: 50%; border: none; padding: 0; cursor: pointer; position: relative; display: flex; align-items: center; justify-content: center; transition: opacity 0.2s; }
.win-btn .icon { opacity: 0; transition: opacity 0.2s; pointer-events: none; }
.window-controls:hover .win-btn .icon { opacity: 0.6; }
.win-btn:hover { opacity: 0.8; }
.win-btn.close { background-color: #ff5f56; border: 0.5px solid #e0443e; }
.win-btn.minimize { background-color: #ffbd2e; border: 0.5px solid #dea123; }
.win-btn.maximize { background-color: #27c93f; border: 0.5px solid #1aab29; }
.win-btn.close .icon::before { content: "×"; color: #4c0000; font-size: 10px; font-weight: bold; }
.win-btn.minimize .icon::before { content: ""; width: 6px; height: 1px; background: #995700; display: block; }
.win-btn.maximize .icon::before { content: ""; width: 6px; height: 6px; border: 1px solid #006500; display: block; clip-path: polygon(0% 0%, 100% 0%, 100% 100%, 0% 100%, 0% 0%, 15% 15%, 15% 85%, 85% 85%, 85% 15%, 15% 15%); }
.container { flex: 1; display: flex; flex-direction: column; padding: 24px 32px 32px; overflow: hidden; }
.main-layout { flex: 1; display: flex; gap: 24px; align-items: stretch; }
.preview-panel { flex: 1; display: flex; flex-direction: column; border: 1px solid var(--border-default); background: var(--surface-subtle); border-radius: 4px; overflow: hidden; min-width: 0; }
.panel-header { padding: 12px 16px; background: rgba(255,255,255,0.02); border-bottom: 1px solid var(--border-default); color: var(--text-muted); font-size: 9px; letter-spacing: 1px; }
.preview-content { flex: 1; display: flex; align-items: center; justify-content: center; padding: 24px; cursor: pointer; transition: background 0.2s; position: relative; }
.preview-content:hover { background: var(--surface-hover); }
.preview-panel.empty .preview-content { background: repeating-linear-gradient(45deg, transparent, transparent 10px, rgba(255,255,255,0.02) 10px, rgba(255,255,255,0.02) 20px); }
.preview-img { max-width: 100%; max-height: 100%; object-fit: contain; filter: drop-shadow(0 4px 24px rgba(0,0,0,0.7)); }
.placeholder { display: flex; flex-direction: column; align-items: center; gap: 12px; pointer-events: none; }
.center-controls { width: 100px; display: flex; flex-direction: column; justify-content: center; gap: 32px; }
.mode-stack { display: flex; flex-direction: column; gap: 8px; }
.mode-btn { background: transparent; border: 1px solid var(--border-default); color: var(--text-muted); padding: 12px 0; font-size: 9px; cursor: pointer; }
.mode-btn.active { background: var(--text-primary); color: var(--bg-color); border-color: var(--text-primary); }
.mode-btn:disabled { opacity: 0.05; cursor: not-allowed; }
.execute-btn { background: transparent; border: 1px solid var(--text-primary); color: var(--text-primary); height: 100px; cursor: pointer; font-size: 13px; }
.execute-btn:hover:not(:disabled) { background: var(--text-primary); color: var(--bg-color); }
.execute-btn:disabled { opacity: 0.05; }
.panel-actions, .panel-actions-placeholder { height: 48px; border-top: 1px solid var(--border-default); background: rgba(0,0,0,0.1); }
.panel-actions { padding: 8px; display: flex; gap: 8px; }
.panel-actions-placeholder { background: transparent; border-top-color: transparent; }
.actions-empty { flex: 1; }
.ghost-btn { flex: 1; background: transparent; border: 1px solid var(--border-strong); color: var(--text-secondary); padding: 0 8px; font-size: 9px; cursor: pointer; border-radius: 2px; height: 100%; }
.ghost-btn:hover { border-color: var(--text-primary); color: var(--text-primary); }
.status-footer { margin-top: 24px; display: flex; justify-content: space-between; align-items: center; border-top: 1px solid var(--border-default); padding-top: 16px; }
.status-wrap.is-error { color: #ff4444; }
.path-info { max-width: 60%; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; opacity: 0.2; font-size: 10px; }
</style>
