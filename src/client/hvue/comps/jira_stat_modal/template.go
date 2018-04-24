package jira_stat_modal

const template = `
<el-dialog 
		:visible.sync="visible" 
		width="60%"
		:before-close="Hide"
>
    <span slot="title" class="novagile">
        <h2 style="margin: 0 0"><i class="fas fa-indent icon--left"></i>Jira Stats</h2>
    </span>

    <hours-tree
            :nodes="nodes"
            @node-click="HandleNodeClick"
    >
    </hours-tree>

</el-dialog>
`
