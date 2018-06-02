package workloadschedule_modal

const template = `
<el-dialog 
		:visible.sync="visible"
		width="70%"
		:before-close="Hide"
>
    <span slot="title" class="prjptf">
        <h2 style="margin: 0 0"><i class="fas fa-chart-bar icon--left"></i>Workload Schedule</h2>
    </span>
    <el-container 	
			v-loading="wrkSchedLoading"
			class="workload-schedule"
	>
        <el-aside>
            <selection-tree
                    ref="selection-tree"
					:wrkSched.sync="wrkSched"
					@update:wrkSched="UpdateBarChart"
            ></selection-tree>
        </el-aside>
        <el-main>
            <div>
                <bars-chart
						v-if="wrkSched"
						ref="bars-chart"
                        :infos="barchartInfos"
                ></bars-chart>
            </div>
        </el-main>
    </el-container>
</el-dialog>
`
