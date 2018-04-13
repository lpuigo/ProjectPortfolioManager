package project_table

/*
    size="mini"
	max-height="100%"
    style="width: 100%"

		resizable="true"
		width="160px"

	style="white-space: pre-wrap;"
*/

const template = `
<el-table
    :data="filteredProjects"
	header-row-class-name="novagile-light"
    :row-class-name="TableRowClassName"
	:default-sort = "{prop: 'client', order: 'ascending'}"
    @current-change="SetSelectedProject"
    @row-dblclick="SelectRow"
	height="100%"
	:border=true
>
    <el-table-column 
            label="Client"	prop="client"	width="160px" sortable :sort-by="['client','name']" 
            :resizable=true :show-overflow-tooltip=true
    ></el-table-column>

    <el-table-column
            label="Project Name"	prop="name"	width="200px"
			:resizable=true :show-overflow-tooltip=true
    ></el-table-column>

    <el-table-column
            label="Comment" prop="comment"
		    :resizable=false
    >
        <template slot-scope="scope">
            <i :class="RiskIconClass(scope.row.risk)"></i><span>{{scope.row.comment}}</span>
        </template>
	</el-table-column>

    <el-table-column 
            label="KickOff"	prop="milestones.Kickoff"	width="100px"	sortable    :sort-by="['milestones.Kickoff', 'client','name']"
		    :resizable=false    align="center"	:formatter="FormatDate"
    ></el-table-column>

    <el-table-column 
            label="UAT"	prop="milestones.UAT"	width="100px"	sortable    :sort-by="['milestones.UAT', 'client','name']"
		    :resizable=false    align="center"	:formatter="FormatDate"
    ></el-table-column>

    <el-table-column 
            label="RollOut"	prop="milestones.RollOut"	width="100px"	sortable    :sort-by="['milestones.RollOut', 'client','name']"
		    :resizable=false    align="center"	:formatter="FormatDate"
    ></el-table-column>

    <el-table-column 
            label="WorkLoad"	width="120px"
		    :resizable=false	align="center"
    >
        <template slot-scope="scope">
            <project-progress-bar :project="scope.row"></project-progress-bar>
        </template>
	</el-table-column>

    <el-table-column
            label="PS"	prop="lead_ps"	width="120px" sortable :sort-by="['lead_ps', 'client','name']"
		    :resizable=false :show-overflow-tooltip=true
 		    :filters="FilterList('lead_ps')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
    ></el-table-column>

    <el-table-column
            label="Lead Dev"	prop="lead_dev"	width="120px" sortable :sort-by="['lead_dev', 'client','name']"
		    :resizable=false :show-overflow-tooltip=true
 		    :filters="FilterList('lead_dev')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
    ></el-table-column>

    <el-table-column
            label="Type"	prop="type"	width="80px"
		    :resizable=false
 		    :filters="FilterList('type')"	:filter-method="FilterHandler"	filter-placement="bottom-end"
    ></el-table-column>

    <el-table-column
            label="Status"	prop="status"	width="120px" sortable :sort-by="['status', 'client','name']"
		    :resizable=false :show-overflow-tooltip=true
 		    :filters="FilterList('status')"	:filter-method="FilterHandler"	filter-placement="bottom-end" :filtered-value="FilteredValue()"
	></el-table-column>
</el-table>
`
