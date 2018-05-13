package workloadschedule_modal

const template = `
<el-dialog 
		:visible.sync="visible"
		width="70%"
		:before-close="Hide"
>
    <span slot="title" class="novagile">
        <h2 style="margin: 0 0"><i class="fas fa-chart-bar icon--left"></i>Workload Schedule</h2>
    </span>

    <div style="min-height: 400px;max-height: 65vh;overflow: auto;">
		<bars-chart
				v-if="wrkSched"
				:weeks="wrkSched.weeks"
				:series="series"
		></bars-chart>
    </div>

</el-dialog>
`
