package wl_progress_bar

const template = `
<div>
    <div v-if="showProgressBar" :class="progressStatus">
        <el-progress 
                 :show-text="false"
                 :stroke-width="5"
                 :percentage="progressPct"
        ></el-progress>
        <span class="small-font">{{ project | chargeFormat }}</span>
    </div>
    <span v-else>{{ project | chargeFormat }}</span>
</div>

`
