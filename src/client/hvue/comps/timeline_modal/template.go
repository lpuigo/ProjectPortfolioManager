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

	<el-container>
		<el-aside width="80px">
			<el-menu :default-active="slotType" class="timeline" @select="HandleSlotType">
				<el-menu-item index="4Q">Year</el-menu-item>
				<el-menu-item index="2Q">Semester</el-menu-item>
				<el-menu-item index="1Q">Quarter</el-menu-item>
			</el-menu>
		</el-aside>
		<el-main style="padding: 0px">
			<div class="timeline">
				<el-table
						:data="timelines"
						:default-sort="{prop: 'milestones', order: 'ascending'}"
						height="100%"
						:border="false"
						size="mini"
						highlight-current-row
				>
					<el-table-column
							label="Project" prop="name" width="240px" :show-overflow-tooltip=true
							sortable
					>
						<template slot-scope="scope">
							<span>{{scope.row.name}}</span>
						</template>
					</el-table-column>

					<el-table-column
							label="Phases" prop="milestones"
							sortable    :sort-by="['milestones.RollOut', 'milestones.GoLive', 'milestones.UAT', 'milestones.Outline', 'milestones.Kickoff', 'name']"
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
		</el-main>
	</el-container>    

</el-dialog>
`
