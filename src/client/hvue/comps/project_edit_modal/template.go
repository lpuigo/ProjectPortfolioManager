package project_edit_modal

const template = `
<el-dialog :visible.sync="visible" width="60%">
    <!--<span slot="title" class="dialog-header">-->
    <span slot="title">
        <h2 v-if="currentProject" style="margin: 0 0"><i class="far fa-edit"></i> Edit Project: <span style="color: teal">{{currentProject.client}} - {{currentProject.name}}</span></h2>
    </span>

	<el-row :gutter="15" class="form-row">
		<el-col :span="12">
			<span><strong>Client</strong></span>
			<el-input placeholder="Client Name" v-model="currentProject.client"></el-input>
		</el-col>	
		<el-col :span="12">
			<span><strong>Project Name</strong></span>
			<el-input placeholder="Project Name" v-model="currentProject.name"></el-input>
		</el-col>	
  	</el-row>

	<el-row :gutter="15" class="form-row">
		<el-col :span="7">
			<span><strong>PS Actor</strong></span>
			<el-input placeholder="PS Name" v-model="currentProject.lead_ps"></el-input>
		</el-col>	
		<el-col :span="7">
			<span><strong>Lead Dev</strong></span>
			<el-input placeholder="Dev Name" v-model="currentProject.lead_dev"></el-input>
		</el-col>	
		<el-col :span="6">
			<span><strong>Project Type</strong></span>
            <el-select v-model="currentProject.type" placeholder="Project Type" filterable class="fluid-select">
                <el-option v-for="vt in typeList" :key="vt.value" :label="vt.text" :value="vt.value"></el-option>
            </el-select>
        </el-col>	
		<el-col :span="4">
			<span><strong>Estim. WL</strong></span>
            <el-input-number v-model="currentProject.forecast_wl"
                             :min="0" 
                             :step="0.5" 
                             controls-position="right"
							 class="fluid-input"
            ></el-input-number>
		</el-col>	
  	</el-row>

	<el-row :gutter="15" class="form-row">
		<el-col :span="24">
			<span><strong>Comment</strong></span>
            <el-input
                    type="textarea" 
                    placeholder="Project comment" 
                    :autosize="{ minRows: 2, maxRows: 8}"
                    v-model="currentProject.comment"
					class="form-textarea"
            ></el-input>
		</el-col>	
  	</el-row>
	
	<el-row :gutter="15" class="form-row">
		<el-col :span="12">
			<span><strong>Status</strong></span>
            <el-select v-model="currentProject.status" placeholder="Project Status" filterable class="fluid-select">
                <el-option v-for="vt in statusList" :key="vt.value" :label="vt.text" :value="vt.value"></el-option>
            </el-select>
		</el-col>	
		<el-col :span="12">
			<span><strong>Risk</strong></span>
            <el-select v-model="currentProject.risk" placeholder="Project Risk" filterable class="fluid-select">
                <el-option v-for="vt in riskList" :key="vt.value" :label="vt.text" :value="vt.value"></el-option>
            </el-select>
		</el-col>	
  	</el-row>
    
	<el-row :gutter="15" class="form-row">
		<el-col :span="24">
            <el-table
                :data="usedMilestoneKeys"
                height="100%"
                :border=true
            >
                <el-table-column 
                        label="Milestone" min-width="120px" 
                        :resizable=false
                >
                    <template slot-scope="scope">
                        <span>{{scope.row}}</span>
                    </template>
                </el-table-column>
                <el-table-column 
                        label="Date" min-width="120px" 
                        :resizable=false
                >
                    <template slot-scope="scope">
                        <span>{{currentProject.milestones[scope.row]}}</span>
                    </template>
                </el-table-column>
            </el-table>
		</el-col>	
  	</el-row>
	
    <span slot="footer" class="dialog-footer">
        <el-popover
                v-if="!isNewProject"
                ref="confirm_delete_popover"
                placement="top"
                width="160"
                v-model="showconfirmdelete"
                :disable="!visible"
        >
            <p>Confirm to delete this project ?</p>
            <div style="text-align: left; margin: 0;">
            	<el-button size="mini" type="text" @click="showconfirmdelete = false">Cancel</el-button>
            	<el-button size="mini" type="primary" @click="DeleteProject">Delete</el-button>
            </div>
        </el-popover>
        
        <el-tooltip effect="light" :open-delay="500">
            <div slot="content">Delete<br/>current project</div>
            <el-button v-if="!isNewProject" type="danger" plain icon="far fa-trash-alt" v-popover:confirm_delete_popover tooltip="Delete"></el-button>
        </el-tooltip>
        <el-tooltip effect="light" :open-delay="500">
            <div slot="content">Create copy of<br/>current project</div>
            <el-button v-if="!isNewProject" type="info" plain icon="far fa-clone" tooltip="Duplicate" @click="Duplicate"></el-button>
        </el-tooltip>
        <el-button @click="visible = false">Cancel</el-button>
        <el-button v-if="!isNewProject" type="success" plain @click="ConfirmChange">Confirm Change</el-button>
        <el-button v-if="isNewProject" type="success" plain @click="NewProject">Create New</el-button>
    </span>
</el-dialog>
`
