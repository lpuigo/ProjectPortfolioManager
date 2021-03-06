package project_stat_modal

const template string = `
<el-dialog 
		:visible.sync="visible" 
		width="60%"
		:before-close="Hide"
>
    <span slot="title" class="prjptf">
        <h2 v-if="project" style="margin: 0 0">
            <i class="fas fa-chart-line icon--left"></i>Project: <span style="color: teal">{{project.client}} - {{project.name}}</span>
        </h2>
    </span>

    <el-tabs 			
			v-model="activeTabName" 
			style="min-height: 300px;"
			@tab-click="ActivateTabs"
	>
        <el-tab-pane label="Issues List" name="issuelist">
			<el-table
					:data="issueInfoList"
					:default-sort = "{prop: 'spent', order: 'descending'}"
					height="60vh"
			>
                <el-table-column type="expand">
                    <template slot-scope="props">
                        <sre-chart 
                                :issuestat="props.row.issueStat"
                                style="height: 150px"
                                :border="true"
                        ></sre-chart>
                    </template>
                </el-table-column>
				<el-table-column 
					label="Issue"	prop="issue"	width="120px"	sortable 
					:resizable=false :show-overflow-tooltip=true
				>
                    <template slot-scope="props">
                        <a :href="props.row.issueStat.href" target="_blank">{{props.row.issue}}</a>
                    </template>
                </el-table-column>
				<el-table-column 
					label="Summary"	prop="summary"	sortable 
					:resizable=false :show-overflow-tooltip=true
				></el-table-column>
				<el-table-column 
					label="Spent"	prop="spent"	width="120px"	sortable 
					:resizable=false :formatter="FormatFloat" align="right"
				></el-table-column>
				<el-table-column 
					label="Remaining"	prop="remaining"	width="120px"	sortable 
					:resizable=false :formatter="FormatFloat" align="right"
				></el-table-column>
				<el-table-column 
					label="% Total Spent"	width="180px" 
					:resizable=false
				>
                    <template slot-scope="props">
                        <el-progress
                                :text-inside="true"
                                :stroke-width="16"
                                :percentage="props.row.projectPct"
                        ></el-progress>
                    </template>
                </el-table-column>
			</el-table>
        </el-tab-pane>
        <el-tab-pane label="Global SRE Chart" name="sre">
            <sre-chart
                    v-if="issueStat"
                    :issuestat="issueStat"
                    style="height: 300px"
                    :border="true"
            ></sre-chart>
        </el-tab-pane>
        <el-tab-pane label="Jira Logs" name="projectlogs">
			<project-tree
					style="max-height: 65vh;overflow: auto;"
					:nodes="plnodes"
					@node-click="HandleNodeClick"
			>
			</project-tree>
        </el-tab-pane>
    </el-tabs>

</el-dialog>
`
