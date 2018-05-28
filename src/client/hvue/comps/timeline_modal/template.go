package timeline_modal

const template string = `
<el-dialog 
		:visible.sync="visible" 
		width="80%"
		:before-close="Hide"
>
    <span slot="title" class="novagile">
        <h2 style="margin: 0 0"><i class="fas fa-stream icon--left"></i>Projects Time Line</h2>
    </span>

    <div class="timeline">
        <el-table
                :data="timelines"
                :default-sort="{prop: 'id', order: 'ascending'}"
                height="100%"
                :border="false"
                size="mini"
        >
            <el-table-column
                    label="Project" prop="id" width="240px" :show-overflow-tooltip=true
                    sortable    :sort-by="['milestones.RollOut', 'milestones.GoLive', 'milestones.UAT', 'milestones.Outline', 'milestones.Kickoff', 'client', 'name']"
            >
                <template slot-scope="scope">
                    <span>{{scope.row.name}}</span>
                </template>
            </el-table-column>

            <el-table-column
                    label="Phases" 
            >
                <template slot-scope="scope">
                    <div class="project-line">
                        <div 
                                v-for="p in scope.row.phases" 
                                class="item" :class="p.name"
                                :style="p.style"
                        ></div>
                    </div>
                </template>
            </el-table-column>
        </el-table>
    </div>

</el-dialog>
`
