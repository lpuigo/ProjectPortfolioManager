package project_table


/*
	max-height="100%"
    style="width: 100%"

		resizable="true"
		width="160px"

	style="white-space: pre-wrap;"
*/

const template = `
<el-table
    :data="filteredProjects"
    :row-class-name="TableRowClassName"
    @current-change="SetSelectedProject"
    @row-dblclick="SelectRow"
	height="100%"
    size="mini"
	:border=true
>
    <el-table-column
        prop="client"	label="Client"	width="160px"
		:resizable=false
    ></el-table-column>
    <el-table-column
        prop="name"	label="Project Name"	width="220px"
    ></el-table-column>
    <el-table-column
        prop="comment"	label="Comment"
		:resizable=false
    ></el-table-column>
    <el-table-column
        prop="lead_ps"	label="PS"	width="120px"
		:resizable=false
    ></el-table-column>
    <el-table-column
        prop="lead_dev"	label="Lead Dev"	width="120px"
		:resizable=false
    ></el-table-column>
    <el-table-column
        prop="type"	label="Type"	width="70px"
		:resizable=false
    ></el-table-column>
    <el-table-column
        prop="status"	label="Status"	width="100px"
		:resizable=false
 		:filters="StatusList()"
		:filter-method="StatusFilter"
		filter-placement="bottom-end"    ></el-table-column>
</el-table>
`
