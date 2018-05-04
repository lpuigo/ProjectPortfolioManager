package project_edit_modal

const template = `
<el-dialog :visible.sync="visible" width="60%">
    <!--<span slot="title" class="dialog-header">-->
    <span slot="title" class="novagile">
        <h2 v-if="currentProject" style="margin: 0 0"><i class="far fa-edit icon--left"></i>Edit Project: <span style="color: teal">{{currentProject.client}} - {{currentProject.name}}</span></h2>
    </span>

	<el-row :gutter="15" class="form-row">
		<el-col :span="12">
			<span><strong>Client</strong></span>
			<el-input placeholder="Client Name" v-model="currentProject.client"></el-input>
		</el-col>	
		<el-col :span="12">
			<span><strong>Project Name</strong></span>
			<el-input placeholder="Project Name" v-model="currentProject.name">
				<el-dropdown slot="prepend"	trigger="click" @command="SetClientName">
					<el-button type="primary" :loading="clientNameLookup" @click="GetClientNameList">
						<i class="fas fa-search icon--right"></i>
					</el-button>
                    <el-dropdown-menu slot="dropdown">
                        <el-dropdown-item v-if="clientNameListEmpty">Retreiving Data ...</el-dropdown-item>
						<el-dropdown-item v-else v-for="(vt, index) in clientNameList" :key="index" :command="vt">{{vt.value}} - {{vt.text}}</el-dropdown-item>
					</el-dropdown-menu>
				</el-dropdown>
			</el-input>
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
                             :step="1" 
                             controls-position="right"
							 class="fluid-input"
                             @change="AuditProject"
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
		<el-col :span="24">
            <el-collapse v-model="displayedInfos" accordion>
                <el-collapse-item name="1">
                    <template slot="title">
                        <span><strong><i class="fas fa-info-circle icon--left"></i>Project warning message(s): {{currentProject.audits.length}}</strong></span>
                    </template>
                    <span v-for="(a, index) in currentProject.audits" :key="index">{{a.priority}} - {{a.title}}</span>
                </el-collapse-item>
            </el-collapse>
		</el-col>	
  	</el-row>
	
	<el-row :gutter="15" class="form-row">
		<el-col :span="12">
			<span><strong>Status</strong></span>
            <el-select v-model="currentProject.status" placeholder="Project Status" filterable class="fluid-select" @change="AuditProject">
                <el-option v-for="vt in statusList" :key="vt.value" :label="vt.text" :value="vt.value"></el-option>
            </el-select>
		</el-col>	
		<el-col :span="12">
			<span><strong>Risk</strong></span>
            <el-select v-model="currentProject.risk" placeholder="Project Risk" filterable class="fluid-select" @change="AuditProject">
                <el-option v-for="vt in riskList" :key="vt.value" :label="vt.text" :value="vt.value"></el-option>
            </el-select>
		</el-col>	
  	</el-row>
    
	<el-row :gutter="15" class="form-row">
        <el-col :span="7">
            <el-dropdown @command="AddMilestone" style="float: right;">
				<el-button type="success" plain size="mini">
					Add a Milestone<i class="el-icon-arrow-down el-icon--right"></i>
				</el-button>
                <el-dropdown-menu slot="dropdown">
                    <el-dropdown-item v-for="ms in unusedMilestoneKeys" :key="ms" :command="ms">{{ms}}</el-dropdown-item>
                </el-dropdown-menu>
            </el-dropdown>
        </el-col>
		<el-col :span="14">
            <el-table
                :data="usedMilestoneKeys"
                :border=true
            >
                <el-table-column 
                        label="Action" width="80px" 
						align="center"	
                        :resizable=false
                >
                    <template slot-scope="scope">
                        <el-button 
                        	type="danger" 
                        	plain
                        	size="mini" 
                        	icon="far fa-calendar-times" 
                        	@click="DeleteMilestone(scope.row)"
                        ></el-button>
                    </template>
                </el-table-column>
                <el-table-column 
                        label="Milestone" width="100px"
						align="right"
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
                        <el-date-picker
                                v-model="currentProject.milestones[scope.row]"
                                type="date"       
                                format="dd/MM/yyyy"
                                value-format="yyyy-MM-dd"
                                :picker-options="{firstDayOfWeek:1}"
                                size="mini"
								:clearable="false"
                                @change="AuditProject"
                        ></el-date-picker>
                    </template>
                </el-table-column>
            </el-table>
		</el-col>
    </el-row>
	
    <span slot="footer" class="dialog-footer">
        <el-popover
                ref="confirm_delete_popover"
                placement="top"
                width="160"
                v-model="showconfirmdelete"
        >
           <p>Confirm to delete this project ?</p>
            <div style="text-align: left; margin: 0;">
            	<el-button size="mini" type="text" @click="showconfirmdelete = false">Cancel</el-button>
            	<el-button size="mini" type="primary" @click="DeleteProject">Delete</el-button>
            </div>
        </el-popover>
        
        <el-tooltip effect="light" :open-delay="500">
            <div slot="content">Delete<br/>current project</div>
            <el-button :disabled="isNewProject" type="danger" plain icon="far fa-trash-alt" v-popover:confirm_delete_popover></el-button>
        </el-tooltip>
        <el-tooltip effect="light" :open-delay="500">
            <div slot="content">Create copy of<br/>current project</div>
            <el-button :disabled="isNewProject" type="info" plain icon="far fa-clone" @click="Duplicate"></el-button>
        </el-tooltip>
        <el-button @click="visible = false">Cancel</el-button>
        <el-button :type="hasWarning" plain @click="ConfirmChange">
        	<span v-if="!isNewProject">Confirm Change</span>
        	<span v-else>Create New</span>
        </el-button>
    </span>
</el-dialog>
`
