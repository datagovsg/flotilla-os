import React, { Component } from "react"
import PropTypes from "prop-types"
import { withRouter } from "react-router-dom"
import { Form as ReactForm } from "react-form"
import { get, isEmpty, omit, has } from "lodash"
import Navigation from "../Navigation/Navigation"
import Loader from "../styled/Loader"
import PopupContext from "../Popup/PopupContext"
import View from "../styled/View"
import Form from "../styled/Form"
import { ReactFormFieldText } from "../Field/FieldText"
import { ReactFormFieldSelect } from "../Field/FieldSelect"
import ReactFormKVField from "../Field/ReactFormKVField"
import TaskContext from "../Task/TaskContext"
import api from "../../api"
import config from "../../config"

import * as requestStateTypes from "../../constants/requestStateTypes"
import intentTypes from "../../constants/intentTypes"

const taskFormTypes = {
  CREATE: "CREATE",
  UPDATE: "UPDATE",
  CLONE: "CLONE",
}

class TaskForm extends Component {
  static transformValues = values =>
    Object.keys(values).reduce((acc, k) => {
      if (k === "memory") {
        acc[k] = +values[k]
      } else {
        acc[k] = values[k]
      }

      return acc
    }, {})

  checkIfAliasExists = alias =>
    new Promise((resolve, reject) => {
      api
        .getTaskByAlias({ alias })
        .then(res => {
          reject({ error: "Alias already exists." })
        })
        .catch(err => {
          resolve({ success: "Good 2 go" })
        })
    })

  handleSubmit = values => {
    const { data, type, push, renderPopup } = this.props

    switch (type) {
      case taskFormTypes.UPDATE:
        api
          .updateTask({
            definitionID: get(data, "definition_id", ""),
            values: TaskForm.transformValues(values),
          })
          .then(responseData => {
            push(`/tasks/${get(responseData, "definition_id", "")}`)
          })
          .catch(this.handleSubmitError)
        break
      case taskFormTypes.CREATE:
      case taskFormTypes.CLONE:
        api
          .createTask({ values })
          .then(responseData => {
            push(`/tasks/${get(responseData, "definition_id", "")}`)
          })
          .catch(this.handleSubmitError)
        break
      default:
        console.warn("TaskForm's `type` prop was not specified, doing nothing.")
    }
  }

  handleSubmitError = error => {
    const { renderPopup } = this.props
    const e = error.getError()

    renderPopup({
      body: e.data,
      intent: intentTypes.error,
      shouldAutohide: false,
      title: "An error occurred",
    })
  }

  renderTitle() {
    switch (this.props.type) {
      case taskFormTypes.CREATE:
        return "Create New Task"
      case taskFormTypes.UPDATE:
        return `Update`
      case taskFormTypes.CLONE:
        return `Clone`
      default:
        return "Task Form"
    }
  }

  shouldNotRenderForm() {
    const { type, requestState } = this.props

    if (
      type !== taskFormTypes.CREATE &&
      requestState === requestStateTypes.NOT_READY
    ) {
      return true
    }

    return false
  }

  getDefaultValues() {
    const { data, type } = this.props

    let ret = {
      memory: get(data, "memory", 1024),
      image: `${get(config, "IMAGE_PREFIX", "")}IMAGE_NAME:IMAGE_TAG`,
    }

    if (type === taskFormTypes.CREATE) {
      return ret
    }

    if (!isEmpty(data)) {
      return {
        ...ret,
        group_name: get(data, "group_name", ""),
        image: get(data, "image", ""),
        command: get(data, "command", ""),
        tags: get(data, "tags", []),
        env: get(data, "env", []),
      }
    }
  }

  getBreadcrumbs = () => {
    const { type, data, definitionID } = this.props

    if (type === taskFormTypes.CREATE) {
      return [
        { text: "Tasks", href: "/tasks" },
        { text: "Create Task", href: "/tasks/create" },
      ]
    }

    const hrefSuffix = type === taskFormTypes.CLONE ? "copy" : "edit"

    return [
      { text: "Tasks", href: "/tasks" },
      {
        text: get(data, "alias", definitionID),
        href: `/tasks/${definitionID}`,
      },
      {
        text: this.renderTitle(),
        href: `/tasks/${definitionID}/${hrefSuffix}`,
      },
    ]
  }

  getActions = ({ shouldDisableSubmitButton }) => {
    const { goBack } = this.props

    return [
      {
        isLink: false,
        text: "Cancel",
        buttonProps: {
          onClick: goBack,
        },
      },
      {
        isLink: false,
        text: "Submit",
        buttonProps: {
          type: "submit",
          intent: intentTypes.primary,
          isDisabled: shouldDisableSubmitButton === true,
        },
      },
    ]
  }

  validateForm = values => {
    let errors = {}

    if (
      this.props.type !== taskFormTypes.UPDATE &&
      (!has(values, "alias") || isEmpty(values.alias))
    ) {
      errors.alias = "Alias cannot be blank."
    }

    if (!has(values, "group_name") || isEmpty(values.group_name)) {
      errors.group_name = "Group name cannot be blank."
    }

    if (!has(values, "image") || isEmpty(values.image)) {
      errors.image = "Image cannot be blank."
    }

    if (!has(values, "command") || isEmpty(values.command)) {
      errors.command = "Command cannot be blank."
    }

    if (!has(values, "memory")) {
      errors.memory = "Memory cannot be blank."
    }

    return errors
  }

  shouldDisableSubmitButton = formAPI => {
    if (!isEmpty(formAPI.errors)) {
      return true
    }

    const requiredValues = ["alias", "group_name", "image", "command", "memory"]

    for (let i = 0; i < requiredValues.length; i++) {
      if (!has(formAPI.values, requiredValues[i])) {
        return true
      }
    }

    return false
  }

  render() {
    const { type } = this.props

    // Don't render the form if, say, the task definition for updating a task
    // has not been fetched. Wait until the next render call.
    if (this.shouldNotRenderForm()) {
      return <Loader />
    }

    return (
      <ReactForm
        defaultValues={this.getDefaultValues()}
        onSubmit={this.handleSubmit}
      >
        {formAPI => {
          console.log(formAPI)
          return (
            <form onSubmit={formAPI.submitForm}>
              <View>
                <Navigation
                  breadcrumbs={this.getBreadcrumbs()}
                  actions={this.getActions({
                    shouldDisableSubmitButton: this.shouldDisableSubmitButton(
                      formAPI
                    ),
                  })}
                />
                <Form title={this.renderTitle()}>
                  {type !== taskFormTypes.UPDATE && (
                    <ReactFormFieldText
                      label="Alias"
                      field="alias"
                      description="Choose a descriptive alias for this task."
                      shouldDebounce
                      isRequired
                      validate={value =>
                        !value ? { error: "Value must not be null." } : null
                      }
                    />
                  )}
                  <ReactFormFieldSelect
                    label="Group Name"
                    field="group_name"
                    requestOptionsFn={api.getGroups}
                    shouldRequestOptions
                    isCreatable
                    isRequired
                    description="Create a new group name or select an existing one to help searching for this task in the future."
                    validate={value =>
                      !value ? { error: "Value must not be null." } : null
                    }
                  />
                  <ReactFormFieldText
                    label="Image"
                    field="image"
                    description="The full URL of the Docker image and tag."
                    isRequired
                    validate={value =>
                      !value ? { error: "Value must not be null." } : null
                    }
                  />
                  <ReactFormFieldText
                    isTextArea
                    label="Command"
                    field="command"
                    description="The command for this task to execute."
                    isRequired
                    validate={value =>
                      !value ? { error: "Value must not be null." } : null
                    }
                  />
                  <ReactFormFieldText
                    isNumber
                    label="Memory (MB)"
                    field="memory"
                    description="The amount of memory this task needs."
                    isRequired
                    validate={value =>
                      !value ? { error: "Value must not be null." } : null
                    }
                  />
                  <ReactFormFieldSelect
                    isCreatable
                    isMulti
                    label="Tags"
                    field="tags"
                    requestOptionsFn={api.getTags}
                    shouldRequestOptions
                    isRequired={false}
                  />
                  <ReactFormKVField
                    label="Environment Variables"
                    field="env"
                    addValue={formAPI.addValue}
                    removeValue={formAPI.removeValue}
                    values={get(formAPI, ["values", "env"], [])}
                    descripion="Environment variables that can be adjusted during execution."
                    isRequired={false}
                    isValueRequired={false}
                    validateKey={value =>
                      !value ? { error: "Key must not be null." } : null
                    }
                    validateValue={() => null}
                  />
                </Form>
              </View>
            </form>
          )
        }}
      </ReactForm>
    )
  }
}

TaskForm.propTypes = {
  data: PropTypes.object,
  goBack: PropTypes.func.isRequired,
  push: PropTypes.func.isRequired,
  renderPopup: PropTypes.func.isRequired,
  requestState: PropTypes.oneOf(Object.values(requestStateTypes)),
  type: PropTypes.oneOf(Object.values(taskFormTypes)).isRequired,
}

const ConnectedTaskForm = withRouter(props => (
  <PopupContext.Consumer>
    {ctx => (
      <TaskForm
        {...omit(props, ["history", "location", "match", "staticContext"])}
        push={props.history.push}
        goBack={props.history.goBack}
        renderPopup={ctx.renderPopup}
      />
    )}
  </PopupContext.Consumer>
))

export const CreateTaskForm = () => (
  <ConnectedTaskForm type={taskFormTypes.CREATE} />
)

export const UpdateTaskForm = props => (
  <TaskContext.Consumer>
    {ctx => <ConnectedTaskForm type={taskFormTypes.UPDATE} {...ctx} />}
  </TaskContext.Consumer>
)

export const CloneTaskForm = props => (
  <TaskContext.Consumer>
    {ctx => <ConnectedTaskForm type={taskFormTypes.CLONE} {...ctx} />}
  </TaskContext.Consumer>
)
