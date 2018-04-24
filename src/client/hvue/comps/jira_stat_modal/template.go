package jira_stat_modal

const template = `
<el-dialog :visible.sync="visible" width="60%">
    <span slot="title" class="novagile">
        <h2 style="margin: 0 0"><i class="far fa-edit icon--left"></i>Jira Stats</h2>
    </span>

    <hours-tree
            :nodes="nodes"
            @node-click="HandleNodeClick"
    >
    </hours-tree>

</el-dialog>
`
