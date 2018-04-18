package project_stat_modal

const template = `
<el-dialog 
		:visible.sync="visible" 
		width="60%"
		:before-close="Hide"
>
    <span slot="title" class="novagile">
        <h2 v-if="project" style="margin: 0 0">
            <i class="fas fa-chart-line icon--left"></i>Project: <span style="color: teal">{{project.client}} - {{project.name}}</span>
        </h2>
    </span>

    <el-tabs tab-position="top" style="min-height: 300px;">
        <el-tab-pane label="Issues List">
			<el-table
					:data="issueInfoList"
					:default-sort = "{prop: 'spent', order: 'descending'}"
					max-height="350"
			>
				<el-table-column 
					label="Issue"	prop="issue"	width="120px"	sortable 
					:resizable=false :show-overflow-tooltip=true
				></el-table-column>
				<el-table-column 
					label="Summary"	prop="summary"	sortable 
					:resizable=false :show-overflow-tooltip=true
				></el-table-column>
				<el-table-column 
					label="Spent"	prop="spent"	width="120px"	sortable 
					:resizable=false :show-overflow-tooltip=true
				></el-table-column>
				<el-table-column 
					label="Remaining"	prop="remaining"	width="120px"	sortable 
					:resizable=false :show-overflow-tooltip=true
				></el-table-column>
			</el-table>
        </el-tab-pane>
        <el-tab-pane label="Global SRE Chart">
			<sre-chart v-if="issueStat"
					:issuestat="issueStat"
					max-height="60vh"
					:border="true"
			></sre-chart>
        </el-tab-pane>

        <el-tab-pane label="Role">Role
        </el-tab-pane>
        <el-tab-pane label="Task">Task
        </el-tab-pane>
    </el-tabs>

</el-dialog>
`
