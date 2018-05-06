package jira_stat_modal

const template = `
<el-dialog 
		:visible.sync="visible" 
		width="70%"
		:before-close="Hide"
>
    <span slot="title" class="novagile">
        <h2 style="margin: 0 0"><i class="fas fa-indent icon--left"></i>Jira Stats</h2>
    </span>

	<el-tabs
			v-model="activeTabName" 
			style="min-height: 300px;"
			@tab-click="ActivateTabs"
	>
		<el-tab-pane label="Weekly Log Summary" name="weeklogs">
			<hours-tree
					:nodes="wlnodes"
					@node-click="HandleNodeClick"
			>
			</hours-tree>
		</el-tab-pane>
		<el-tab-pane label="Project Log Summary" name="projectlogs">
			<project-tree
					style="max-height: 65vh;overflow: auto;"
					:nodes="plnodes"
					@node-click="HandleNodeClick"
			>
			</hours-tree>
		</el-tab-pane>
	</el-tabs>

</el-dialog>
`
