package project_table

const template = `
<el-table
    :data="projects"
    height="auto"
    style="width: 100%"
    :row-class-name="TableRowClassName"
    @current-change="SetSelectedProject"
    @row-dblclick="SelectRow"
    size="mini"
>
    <el-table-column
        prop="name"
        label="Project Name"
    ></el-table-column>
    <el-table-column
        prop="status"
        label="Status"
    >
    </el-table-column>
    <el-table-column
        prop="workload"
        label="Estimated Work Load"
    >
    </el-table-column>
</el-table>
`
