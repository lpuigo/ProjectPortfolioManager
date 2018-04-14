package project_edit_modal

const template = `
<el-dialog :visible.sync="visible" width="60%">
    <!--<span slot="title" class="dialog-header">-->
    <span slot="title">
        <h2 v-if="edited_project" style="margin: 0 0">Edit Project: <span style="color: teal">{{edited_project.client}} - {{edited_project.name}}</span></h2>
    </span>
        
    <el-form v-if="currentProject"
        :model="currentProject"
        label-width="150px"
        size="mini"
        >
        <el-form-item label="Client">
            <el-input v-model="currentProject.client"></el-input>
        </el-form-item>
        <el-form-item label="Project Name">
            <el-input v-model="currentProject.name"></el-input>
        </el-form-item>
        <el-form-item label="Project Risk">
            <el-select v-model="currentProject.risk" placeholder="Risk level">
                <el-option v-for="vt in riskList" :label="vt.text" :value="vt.value"></el-option>
            </el-select>
        </el-form-item>
        <el-form-item label="Project Status">
            <el-select v-model="currentProject.status" placeholder="Project status">
                <el-option v-for="vt in statusList" :label="vt.text" :value="vt.value"></el-option>
            </el-select>
        </el-form-item>
        <el-form-item label="Project Type">
            <el-select v-model="currentProject.type" placeholder="Project Type">
                <el-option v-for="vt in typeList" :label="vt.text" :value="vt.value"></el-option>
            </el-select>
        </el-form-item>
        <el-form-item label="Forecast WorkLoad">
            <el-input-number v-model="currentProject.forecast_wl"
                    :min="0"
                    :step="0.5"
            ></el-input-number>
        </el-form-item>
        <el-form-item label="Comment">
            <el-input
                    type="textarea" placeholder="Project comment" :autosize="{ minRows: 2, maxRows: 6}"
                    v-model="currentProject.comment"
            ></el-input>
        </el-form-item>
    </el-form>
    
    <span slot="footer" class="dialog-footer">
        <el-popover
                v-if="!isNewProject"
                ref="popoverdelete"
                placement="top"
                width="160"
                v-model="showconfirmdelete"
                :disable="!visible"
        >
            <p>Confirm to delete this project ?</p>
            <div style="text-align: right; margin: 0">
            <el-button size="mini" type="text" @click="showconfirmdelete = false">Cancel</el-button>
            <el-button size="mini" type="primary" @click="RemoveProject">Delete</el-button>
            </div>
        </el-popover>
        
        <el-button v-if="!isNewProject" type="danger" icon="el-icon-delete" v-popover:popoverdelete></el-button>
        <el-button v-if="!isNewProject" @click="Duplicate">Duplicate</el-button>
        <el-button @click="visible = false">Cancel</el-button>
        <el-button v-if="!isNewProject" type="primary" @click="ConfirmChange">Confirm Change</el-button>
        <el-button v-if="isNewProject" type="primary" @click="NewProject">Create New</el-button>
    </span>
</el-dialog>
`
