@include("components.buttons.template",[
    "icon"=>$icon??"edit",
    "color"=>$color??"primary",
    "class"=>$class??"",
    "message"=>$message??null,
    "function"=>$function??"showFormModal($id)",
    "label"=>$label??"Edit",
    "toggleModal"=>$toggleModal??"#dataModal"
])